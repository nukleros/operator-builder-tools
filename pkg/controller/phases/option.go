// SPDX-License-Identifier: MIT

package phases

import ctrl "sigs.k8s.io/controller-runtime"

type PhaseOption func(*Phase)

func WithCustomRequeueResult(requeueResult ctrl.Result) PhaseOption {
	return func(p *Phase) {
		p.requeueResult = &requeueResult
	}
}
