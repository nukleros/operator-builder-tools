// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

// CheckReadyPhase executes checking for a parent component's readiness status.
func CheckReadyPhase(r workload.Reconciler, req *workload.Request, options ...ResourceOption) (bool, error) {
	// check to see if known types are ready
	knownReady, err := resourcesAreReady(r, req)
	if err != nil {
		return false, fmt.Errorf("unable to determine if resources are ready, %w", err)
	}

	// check to see if the custom methods return ready
	customReady, err := r.CheckReady(req)
	if err != nil {
		return false, fmt.Errorf("unable to determine if resources are ready, %w", err)
	}

	return (knownReady && customReady), nil
}

// resourcesAreReady gets the resources in memory, pulls the current state from the
// clusters and determines if they are in a ready condition.
func resourcesAreReady(r workload.Reconciler, req *workload.Request) (bool, error) {
	// get resources in memory
	desiredResources, err := r.GetResources(req)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve resources, %w", err)
	}

	// get resources from cluster
	clusterResources := make([]client.Object, len(desiredResources))

	for i, rsrc := range desiredResources {
		clusterResource, err := resources.Get(r, req, rsrc)
		if err != nil {
			return false, fmt.Errorf("unable to retrieve resource %s, %w", rsrc.GetName(), err)
		}

		clusterResources[i] = clusterResource
	}

	// check to see if known types are ready
	return resources.AreReady(clusterResources...)
}
