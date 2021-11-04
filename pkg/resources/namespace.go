/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NamespaceKind    = "Namespace"
	NamespaceVersion = "v1"
)

type NamespaceResource struct {
	Object v1.Namespace
}

// NewNamespaceResource creates and returns a new NamespaceResource.
func NewNamespaceResource(object metav1.Object) (*NamespaceResource, error) {
	namespace := &v1.Namespace{}

	err := ToProper(namespace, object)
	if err != nil {
		return nil, err
	}

	return &NamespaceResource{Object: *namespace}, nil
}

// IsReady defines the criteria for a namespace to be condsidered
// ready.
func (namespace *NamespaceResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if namespace.Object.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.Object.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.Object.Status.Phase == v1.NamespaceActive, nil
}
