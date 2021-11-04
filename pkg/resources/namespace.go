/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NamespaceKind    = "Namespace"
	NamespaceVersion = "v1"
)

type NamespaceResource struct {
	v1.Namespace
}

// NewNamespaceResource creates and returns a new NamespaceResource.
func NewNamespaceResource(name, namespace string) *NamespaceResource {
	return &NamespaceResource{
		v1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       NamespaceKind,
				APIVersion: NamespaceVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// GetParent returns the parent attribute of the resource.
func (namespace *NamespaceResource) GetParent() client.Object {
	return namespace
}

// IsReady defines the criteria for a namespace to be condsidered
// ready.
func (namespace *NamespaceResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if namespace.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.Status.Phase == v1.NamespaceActive, nil
}
