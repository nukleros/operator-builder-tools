/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NamespaceKind = "Namespace"
)

type NamespaceResource struct {
	parent *v1.Namespace
}

// NewNamespaceResource creates and returns a new NamespaceResource.
func NewNamespaceResource() *NamespaceResource {
	return &NamespaceResource{
		parent: &v1.Namespace{},
	}
}

// GetParent returns the parent attribute of the resource.
func (namespace *NamespaceResource) GetParent() client.Object {
	return namespace.parent
}

// IsReady defines the criteria for a namespace to be condsidered
// ready.
func (namespace *NamespaceResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if namespace.parent.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.parent.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.parent.Status.Phase == v1.NamespaceActive, nil
}

// NamespaceForResourceIsReady checks to see if the namespace of a resource is
// ready.
func NamespaceForResourceIsReady(rsrc *Resource) (bool, error) {
	// create a stub namespace resource to pass to the NamespaceIsReady method
	namespace := &Resource{
		Client: rsrc.Client,
	}

	// insert the inherited fields
	namespace.Name = rsrc.Namespace
	namespace.Group = ""
	namespace.Version = "v1"
	namespace.Kind = NamespaceKind

	rsrc.setResourceChecker()

	// get the object from the kubernetes cluster
	if err := GetObject(rsrc, true); err != nil {
		return false, err
	}

	return rsrc.resourceChecker.IsReady(rsrc)
}
