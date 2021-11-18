// SPDX-License-Identifier: MIT

package predicates

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

// ResourcePredicates returns the filters which are used to filter out the common reconcile events
// prior to reconciling the child resource of a component.
func ResourcePredicates(r workload.Reconciler, req *workload.Request) predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return needsReconciliation(
				r,
				req,
				e.ObjectOld,
				e.ObjectNew,
			)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			// do not run reconciliation on unknown events
			return false
		},
		CreateFunc: func(e event.CreateEvent) bool {
			// do not run reconciliation again when we just created the child resource
			return false
		},
	}
}

// WorkloadPredicates returns the filters which are used to filter out the common reconcile events
// prior to reconciling an object for a component.
func WorkloadPredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

// needsReconciliation performs some simple checks and returns whether or not a
// resource needs to be updated.
func needsReconciliation(r workload.Reconciler, req *workload.Request, existing, requested client.Object) bool {
	// skip if the objects support observed generation and they are equal
	if existing.GetGeneration() > 0 && requested.GetGeneration() > 0 {
		if existing.GetGeneration() != requested.GetGeneration() {
			return true
		}
	}

	if existing.GetGeneration() == 0 && requested.GetGeneration() == 0 {
		return true
	}

	// get the desired object from the reconciler and ensure that we both
	// found that desired object and that the desired object fields are equal
	// to the existing object fields
	desired, err := GetDesiredObject(r, req, requested)
	if err != nil {
		r.GetLogger().Error(
			err, "unable to get object in memory",
			"kind", requested.GetObjectKind().GroupVersionKind().Kind,
			"name", requested.GetName(),
			"namespace", requested.GetNamespace(),
		)

		return false
	}

	if desired == nil {
		return true
	}

	equal, err := resources.AreEqual(desired, requested)
	if err != nil {
		r.GetLogger().Error(
			err, "unable to determine equality for reconciliation",
			"kind", desired.GetObjectKind().GroupVersionKind().Kind,
			"name", desired.GetName(),
			"namespace", desired.GetNamespace(),
		)

		return true
	}

	return !equal
}

// GetDesiredObject returns the desired object from a list stored in memory.
func GetDesiredObject(r workload.Reconciler, req *workload.Request, compared client.Object) (client.Object, error) {
	var desired client.Object

	allObjects, err := r.GetResources(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get resources, %w", err)
	}

	for _, resource := range allObjects {
		if resources.EqualGVK(compared, resource) && resources.EqualNamespaceName(compared, resource) {
			return resource, nil
		}
	}

	return desired, nil
}
