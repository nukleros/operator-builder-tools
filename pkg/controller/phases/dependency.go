// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

// DependencyPhase executes a dependency check prior to attempting to create resources.
func DependencyPhase(r workload.Reconciler, req *workload.Request, options ...ResourceOption) (bool, error) {
	if !req.Workload.GetDependencyStatus() {
		satisfied, err := dependenciesSatisfied(r, req)
		if err != nil {
			return false, fmt.Errorf("unable to list dependencies, %w", err)
		}

		return satisfied, nil
	}

	return true, nil
}

// dependenciesSatisfied will return whether or not all dependencies are satisfied for a component.
func dependenciesSatisfied(r workload.Reconciler, req *workload.Request) (bool, error) {
	for _, dep := range req.Workload.GetDependencies() {
		satisfied, err := dependencySatisfied(r, req, dep)
		if err != nil || !satisfied {
			return false, err
		}
	}

	return true, nil
}

// dependencySatisfied will return whether or not an individual dependency is satisfied.
func dependencySatisfied(r workload.Reconciler, req *workload.Request, dependency workload.Workload) (bool, error) {
	// get the dependencies by kind that already exist in cluster
	dependencyList := &unstructured.UnstructuredList{}

	dependencyList.SetGroupVersionKind(dependency.GetWorkloadGVK())

	if err := r.List(req.Context, dependencyList, &client.ListOptions{}); err != nil {
		return false, fmt.Errorf("unable to list dependencies, %w", err)
	}

	// expect only one item returned, otherwise dependencies are considered unsatisfied
	if len(dependencyList.Items) != 1 {
		return false, nil
	}

	// get the status.created field on the object and return the status and any errors found
	status, found, err := unstructured.NestedBool(dependencyList.Items[0].Object, "status", "created")
	if err != nil {
		return false, fmt.Errorf("unable to retrieve status.created field, %w", err)
	}

	if !found {
		return false, nil
	}

	return status, nil
}
