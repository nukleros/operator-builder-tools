// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/reconcile"
	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/resources"
	"github.com/nukleros/operator-builder-tools/pkg/status"
)

// CreateResourcesPhase creates or updated the child resources of a workload during a reconciliation loop.
func CreateResourcesPhase(r workload.Reconciler, req *workload.Request, options ...ResourceOption) (bool, error) {
	// get the resources in memory
	desiredResources, err := r.GetResources(req)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve resources, %w", err)
	}

	proceed := true

	for _, resource := range desiredResources {
		condition, created, err := HandleResourcePhaseExit(
			persistResourcePhase(r, req, resource),
		)
		if err != nil {
			if !IsOptimisticLockError(err) {
				req.Log.Error(err, "unable to create or update resource")
			}
		}

		resourceObject := status.ToCommonResource(resource)
		resourceObject.ChildResourceCondition = condition

		// update the status conditions and return any errors
		if err := UpdateResourceConditions(r, req, resourceObject); err != nil {
			if !IsOptimisticLockError(err) {
				r.GetLogger().Error(
					err, "failed to update resource conditions",
					"kind", resource.GetObjectKind().GroupVersionKind().Kind,
					"name", resource.GetName(),
					"namespace", resource.GetNamespace(),
				)

				created = false
			}
		}

		proceed = proceed && created
	}

	return proceed, err
}

// UpdateResourceConditions updates the status.resourceConditions field of the parent custom resource.
func UpdateResourceConditions(
	r workload.Reconciler,
	req *workload.Request,
	resource *status.ChildResource,
) error {
	req.Workload.SetChildResourceCondition(resource)

	if err := r.Status().Update(req.Context, req.Workload); err != nil {
		return fmt.Errorf("unable to update Resource Condition for %s, %w", req.Workload.GetWorkloadGVK().Kind, err)
	}

	return nil
}

// HandleResourcePhaseExit will generate the appropriate resource condition for a resource creation event.
func HandleResourcePhaseExit(
	resourceCreated bool,
	resourceErr error,
) (status.ChildResourceCondition, bool, error) {
	if resourceErr != nil {
		if !IsOptimisticLockError(resourceErr) {
			return status.GetFailResourceCondition(resourceErr), false, resourceErr
		}
	}

	if !resourceCreated {
		return status.GetPendingResourceCondition(), resourceCreated, nil
	}

	return status.GetSuccessResourceCondition(), true, nil
}

// persistResourcePhase executes persisting resources to the Kubernetes database.
func persistResourcePhase(
	r workload.Reconciler,
	req *workload.Request,
	resource client.Object,
) (bool, error) {
	ready, err := commonWait(r, req, resource)
	if err != nil {
		return false, err
	}

	// return the result if the object is not ready
	if !ready {
		return false, nil
	}

	// persist the resource
	if err := CreateOrUpdate(r, req, resource); err != nil {
		if IsOptimisticLockError(err) {
			return true, nil
		}

		return false, fmt.Errorf("unable to create or update resource %s, %w", resource.GetName(), err)
	}

	return true, nil
}

// CreateOrUpdate creates a resource if it does not already exist or updates a resource
// if it does already exist.
func CreateOrUpdate(r workload.Reconciler, req *workload.Request, resource client.Object) error {
	// set ownership on the underlying resource being created or updated
	if err := ctrl.SetControllerReference(req.Workload, resource, r.Scheme()); err != nil {
		req.Log.Error(
			err, "unable to set owner reference on resource",
			"resourceName", resource.GetName(),
			"resourceNamespace", resource.GetNamespace(),
		)

		return fmt.Errorf("unable to set owner reference on %s, %w", resource.GetName(), err)
	}

	// get the resource from the cluster
	clusterResource, err := resources.Get(r, req, resource)
	if err != nil {
		return fmt.Errorf("unable to retrieve resource %s, %w", resource.GetName(), err)
	}

	// create the resource if we have a nil object, or update the resource if we have one
	// that exists in the cluster already
	if clusterResource == nil {
		if err := resources.Create(r, req, resource); err != nil {
			return fmt.Errorf("unable to create resource %s, %w", resource.GetName(), err)
		}

		// add the created event
		status.Created.RegisterAction(r.GetEventRecorder(), resource, req.Workload)

		return reconcile.Watch(r, req, resource)
	}

	// add the updated event
	status.Updated.RegisterAction(r.GetEventRecorder(), resource, req.Workload)

	if err := resources.Update(r, req, resource, clusterResource); err != nil {
		return fmt.Errorf("unable to update resource %s, %w", resource.GetName(), err)
	}

	return nil
}
