// SPDX-License-Identifier: MIT

package phases

import ctrl "sigs.k8s.io/controller-runtime"

// PhaseOption is a function pattern to allow customization of a phase upon registration.
type PhaseOption func(*Phase)

// ResourceOption is a pattern to allow customization to the resource deployment process.
type ResourceOption int

const (
	ResorceOptionWithWait = iota
)

// WithCustomRequeueResult allows you to define a custom result for a phase when it is requeued,
// this allows for custom requeue time backoffs.
func WithCustomRequeueResult(requeueResult ctrl.Result) PhaseOption {
	return func(p *Phase) {
		p.requeueResult = &requeueResult
	}
}

// WithWait adds the ResourceOptionWithWait resource option to the phase.
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
