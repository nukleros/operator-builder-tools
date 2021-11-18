// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/status"
)

// HandlerFunc is an adapter ot allow the use of ordinary functions as reconcile phases.
// if the function has the appropriate signature, it will considered a valid phase handler.
type HandlerFunc func(r workload.Reconciler, req *workload.Request) (proceed bool, err error)

// Phase defines a phase of the reconciliation process.
type Phase struct {
	Name       string
	definition HandlerFunc
}

// DefaultRequeue executes checking for a parent components readiness status.
func (*Phase) DefaultRequeue() ctrl.Result {
	return ctrl.Result{
		Requeue: true,
		// RequeueAfter: 5 * time.Second,
	}
}

// HandlePhaseExit will perform the steps required to exit a phase.
func (p *Phase) HandlePhaseExit(
	r client.StatusClient,
	req *workload.Request,
	phaseIsReady bool,
	phaseError error,
) (ctrl.Result, error) {
	var condition status.PhaseCondition

	var result ctrl.Result

	switch {
	case phaseError != nil:
		if IsOptimisticLockError(phaseError) {
			phaseError = nil
			condition = status.GetSuccessCondition(p.Name)
		} else {
			condition = status.GetFailCondition(p.Name, phaseError)
		}

		result = DefaultReconcileResult()
	case !phaseIsReady:
		condition = status.GetPendingCondition(p.Name)
		result = p.DefaultRequeue()
	default:
		condition = status.GetSuccessCondition(p.Name)
		result = DefaultReconcileResult()
	}

	// update the status conditions and return any errors
	if updateError := updatePhaseConditions(r, req, &condition); updateError != nil {
		// adjust the message if we had both an update error and a phase error
		if !IsOptimisticLockError(updateError) {
			if phaseError != nil {
				phaseError = fmt.Errorf("failed to update status conditions; %v; %w", updateError, phaseError)
			} else {
				phaseError = updateError
			}
		}
	}

	return result, phaseError
}

// updatePhaseConditions updates the status.conditions field of the parent custom resource.
func updatePhaseConditions(r client.StatusClient, req *workload.Request, condition *status.PhaseCondition) error {
	req.Workload.SetPhaseCondition(condition)

	if err := r.Status().Update(req.Context, req.Workload); err != nil {
		return fmt.Errorf("unable to update Phase Condition for %s, %w", req.Workload.GetWorkloadGVK().Kind, err)
	}

	return nil
}
