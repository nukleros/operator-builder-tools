// SPDX-License-Identifier: MIT

package phases

import (
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
)

const optimisticLockErrorMsg = "the object has been modified; please apply your changes to the latest version and try again"

// Requeue will return the default result to requeue a reconciler request when needed.
func Requeue() ctrl.Result {
	return ctrl.Result{Requeue: true}
}

// IsOptimisticLockError checks to see if the error is a locking error.
func IsOptimisticLockError(err error) bool {
	return strings.Contains(err.Error(), optimisticLockErrorMsg)
}

// DefaultReconcileResult will return the default reconcile result when requeuing is not needed.
func DefaultReconcileResult() ctrl.Result {
	return ctrl.Result{}
}
