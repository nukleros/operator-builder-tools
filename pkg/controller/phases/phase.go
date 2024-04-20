// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/status"
)

// HandlerFunc is an adapter to allow the use of ordinary functions as reconcile phases.
// If the function has the appropriate signature, it will considered a valid phase handler.
type HandlerFunc func(r workload.Reconciler, req *workload.Request, options ...ResourceOption) (proceed bool, err error)

// Phase defines a phase of the reconciliation process.
type Phase struct {
	Name            string
	definition      HandlerFunc
	requeueResult   *ctrl.Result
	resourceOptions []ResourceOption
}

// Requeue will return the phase's reconcile result when requeueing is needed.
func (p *Phase) Requeue() ctrl.Result {
	if p.requeueResult != nil {
		return *p.requeueResult
	}

	return ctrl.Result{
		Requeue: true,
	}
}

// DefaultReconcileResult will return the default reconcile result when requeuing is not needed.
func (p *Phase) DefaultReconcileResult() ctrl.Result {
	return ctrl.Result{}
}

// handlePhaseExit will perform the steps required to exit a phase.
func (p *Phase) handlePhaseExit(
	r workload.Reconciler,
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

		result = ctrl.Result{}
	case !phaseIsReady:
		condition = status.GetPendingCondition(p.Name)
		result = p.Requeue()
	default:
		condition = status.GetSuccessCondition(p.Name)
		result = p.DefaultReconcileResult()
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
func updatePhaseConditions(r workload.Reconciler, req *workload.Request, condition *status.PhaseCondition) error {
	req.Workload.SetPhaseCondition(condition)

	if err := r.Status().Update(req.Context, req.Workload); err != nil {
		return fmt.Errorf("unable to update Phase Condition for %s, %w", req.Workload.GetWorkloadGVK().Kind, err)
	}

	return nil
}
