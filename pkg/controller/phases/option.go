// SPDX-License-Identifier: MIT

package phases

import ctrl "sigs.k8s.io/controller-runtime"

// PhaseOption is a function pattern to allow customization of a phase upon registration.
type PhaseOption func(*Phase)

// ResourceOption is a pattern to allow customization to the resource deployment process.
type ResourceOption int

const (
	ResourceOptionWithWait = iota
)

// WithCustomRequeueResult allows you to define a custom result for a phase when it is requeued,
// this allows for custom requeue time backoffs.
func WithCustomRequeueResult(requeueResult ctrl.Result) PhaseOption {
	return func(p *Phase) {
		p.requeueResult = &requeueResult
	}
}

// WithResourceOptions adds the requested resource options to the phase.
func WithResourceOptions(options ...ResourceOption) PhaseOption {
	return func(p *Phase) {
		if len(options) == 0 {
			return
		}

		if len(p.resourceOptions) == 0 {
			p.resourceOptions = options

			return
		}

		p.resourceOptions = append(p.resourceOptions, options...)
	}
}

// hasResourceOption returns true if a set of resource options has a given
// resource option.
func hasResourceOption(option ResourceOption, options ...ResourceOption) bool {
	for i := range options {
		if options[i] == option {
			return true
		}
	}

	return false
}
