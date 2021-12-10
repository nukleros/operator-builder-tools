// SPDX-License-Identifier: MIT

package phases

import ctrl "sigs.k8s.io/controller-runtime"

// PhaseOption is a function pattern to allow customization of a phase upon registration.
type PhaseOption func(*Phase)

// WithCustomRequeueResult allows you to define a custom result for a phase when it is requeued,
// this allows for custom requeue time backoffs.
func WithCustomRequeueResult(requeueResult ctrl.Result) PhaseOption {
	return func(p *Phase) {
		p.requeueResult = &requeueResult
	}
}
