// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

// TODO: the following allows all controllers to list all namespaces,
// regardless of whether or not the controller manages namespaces.
//
// This will eventually be moved into a validating webhook so that the user
// will get a message outlining their mistake rather than buried in the
// reconciliation loop, causing pain when having to sift through logs to
// determine a problem.
//
// See:
//   - https://github.com/vmware-tanzu-labs/operator-builder/issues/141
//   - https://github.com/vmware-tanzu-labs/operator-builder/issues/162

// commonWait applies all common waiting functions for known resources.
func commonWait(
	r workload.Reconciler,
	req *workload.Request,
	resource client.Object,
) (bool, error) {
	// Namespace
	if resource.GetNamespace() != "" {
		ready, err := resources.NamespaceForResourceIsReady(r, req, resource)
		if err != nil {
			return ready, fmt.Errorf("unable to determine if %s namespace is ready, %w", resource.GetNamespace(), err)
		}

		return ready, nil
	}

	return true, nil
}
