// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

// DeletionCompletePhase executes the completion of a reconciliation loop for a delete request.
func DeletionCompletePhase(r workload.Reconciler, req *workload.Request) (bool, error) {
	req.Log.Info("successfully deleted")

	return true, nil
}

func RegisterDeleteHooks(r workload.Reconciler, req *workload.Request) error {
	myFinalizerName := fmt.Sprintf("%s/Finalizer", req.Workload.GetWorkloadGVK().Group)

	if req.Workload.GetDeletionTimestamp().IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(req.Workload.GetFinalizers(), myFinalizerName) {
			controllerutil.AddFinalizer(req.Workload, myFinalizerName)

			if err := r.Update(req.Context, req.Workload); err != nil {
				return fmt.Errorf("unable to register delete hook on %s, %w", req.Workload.GetWorkloadGVK().Kind, err)
			}
		}
	}

	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}
