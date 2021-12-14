// SPDX-License-Identifier: MIT

package reconcile

import (
	"fmt"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/nukleros/operator-builder-tools/pkg/controller/predicates"
	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

// Watch watches a resource.
func Watch(
	r workload.Reconciler,
	req *workload.Request,
	resource client.Object,
) error {
	// ignore jobs as they are ephemerial
	if resource.GetObjectKind().GroupVersionKind().Kind == "job" {
		if resource.GetObjectKind().GroupVersionKind().Version == "v1" {
			return nil
		}
	}

	// check if the resource is already being watched
	var watched bool

	if len(r.GetWatches()) > 0 {
		for _, watcher := range r.GetWatches() {
			if reflect.DeepEqual(
				resource.GetObjectKind().GroupVersionKind(),
				watcher.GetObjectKind().GroupVersionKind(),
			) {
				watched = true
			}
		}
	}

	// watch the resource if it current is not being watched
	if !watched {
		if err := r.GetController().Watch(
			&source.Kind{Type: resource},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    req.Workload,
			},
			predicates.ResourcePredicates(r, req),
		); err != nil {
			return fmt.Errorf("unable to watch resource, %w", err)
		}

		r.SetWatch(resource)
	}

	return nil
}
