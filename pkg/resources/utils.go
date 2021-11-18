// SPDX-License-Identifier: MIT

package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

// Create creates a resource.
func Create(r workload.Reconciler, req *workload.Request, resource client.Object) error {
	r.GetLogger().Info(
		"creating resource",
		"kind", resource.GetObjectKind().GroupVersionKind().Kind,
		"name", resource.GetName(),
		"namespace", resource.GetNamespace(),
	)

	if err := r.Create(
		req.Context,
		resource,
		&client.CreateOptions{FieldManager: r.GetFieldManager()},
	); err != nil {
		return fmt.Errorf("unable to create resource; %w", err)
	}

	return nil
}

// Get gets a resource.
func Get(r workload.Reconciler, req *workload.Request, resource client.Object) (client.Object, error) {
	// create a stub object to store the current resource in the cluster so that we do not affect
	// the desired state of the resource object in memory
	resourceStore := &unstructured.Unstructured{}
	resourceStore.SetGroupVersionKind(resource.GetObjectKind().GroupVersionKind())

	if err := r.Get(
		req.Context,
		client.ObjectKeyFromObject(resource),
		resourceStore,
	); err != nil {
		if errors.IsNotFound(err) {
			// return nil here so we can easily determine if a resource was not
			// found without having to worry about its type
			return nil, nil
		}

		return nil, fmt.Errorf("unable to get resource %s, %w", resource.GetName(), err)
	}

	return resourceStore, nil
}

// Update updates a resource.
func Update(r workload.Reconciler, req *workload.Request, newResource, oldResource client.Object) error {
	needsUpdate, err := NeedsUpdate(r, newResource, oldResource)
	if err != nil {
		return err
	}

	if needsUpdate {
		r.GetLogger().Info(
			"updating resource",
			"kind", oldResource.GetObjectKind().GroupVersionKind().Kind,
			"name", oldResource.GetName(),
			"namespace", oldResource.GetNamespace(),
		)

		if err := r.Patch(
			req.Context,
			newResource,
			client.Merge,
			&client.PatchOptions{FieldManager: r.GetFieldManager()},
		); err != nil {
			return fmt.Errorf("unable to update resource; %w", err)
		}
	}

	return nil
}

// NeedsUpdate determines if a resource needs to be updated.
func NeedsUpdate(r workload.Reconciler, desired, actual client.Object) (bool, error) {
	// check for equality first as this will let us avoid spamming user logs
	// when resources that need to be skipped explicitly (e.g. CRDs) are seen
	// as equal anyway
	equal, err := AreEqual(desired, actual)
	if err != nil {
		return false, fmt.Errorf("unable to determine if resources are equal, %w", err)
	}

	if equal {
		return false, nil
	}

	// always skip custom resource updates as they are sensitive to modification
	// e.g. resources provisioned by the resource definition would not
	// understand the update to a spec
	if desired.GetObjectKind().GroupVersionKind().Kind == "CustomResourceDefinition" {
		messageVerbose := fmt.Sprintf("if updates are desired, consider re-deploying the parent " +
			"resource or generating a new api version with the desired " +
			"changes")

		r.GetLogger().V(4).Info("skipping update", "CustomResourceDefinition", desired.GetName())
		r.GetLogger().V(7).Info(messageVerbose, "CustomResourceDefinition", desired.GetName())

		return false, nil
	}

	return true, nil
}
