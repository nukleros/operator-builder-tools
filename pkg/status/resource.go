// SPDX-License-Identifier: MIT

package status

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ChildResource is the resource and its condition as stored on the workload custom resource's status field.
type ChildResource struct {
	// Group defines the API Group of the resource.
	Group string `json:"group"`

	// Version defines the API Version of the resource.
	Version string `json:"version"`

	// Kind defines the kind of the resource.
	Kind string `json:"kind"`

	// Name defines the name of the resource from the metadata.name field.
	Name string `json:"name"`

	// Namespace defines the namespace in which this resource exists in.
	Namespace string `json:"namespace"`

	// ResourceCondition defines the current condition of this resource.
	ChildResourceCondition `json:"condition,omitempty"`
}

// ChildResourceCondition describes the condition of a Kubernetes resource managed by the parent object.
type ChildResourceCondition struct {
	// Created defines whether this object has been successfully created or not.
	Created bool `json:"created"`

	// LastModified defines the time in which this resource was updated.
	LastModified string `json:"lastModified,omitempty"`

	// Message defines a helpful message from the resource phase.
	Message string `json:"message,omitempty"`
}

// ToCommonResource converts a client.Object into a common API resource.
func ToCommonResource(resource client.Object) *ChildResource {
	resourceCommon := &ChildResource{
		Group:     resource.GetObjectKind().GroupVersionKind().Group,
		Version:   resource.GetObjectKind().GroupVersionKind().Version,
		Kind:      resource.GetObjectKind().GroupVersionKind().Kind,
		Name:      resource.GetName(),
		Namespace: resource.GetNamespace(),
	}

	return resourceCommon
}

// GetSuccessResourceCondition defines the success condition for the phase.
func GetSuccessResourceCondition() ChildResourceCondition {
	return ChildResourceCondition{
		Created:      true,
		LastModified: time.Now().UTC().String(),
		Message:      "resource creation successful",
	}
}

// GetPendingResourceCondition defines the pending condition for the phase.
func GetPendingResourceCondition() ChildResourceCondition {
	return ChildResourceCondition{
		Created:      false,
		LastModified: time.Now().UTC().String(),
		Message:      "Pending Execution of Resource Creation",
	}
}

// GetFailResourceCondition defines the fail condition for the phase.
func GetFailResourceCondition(err error) ChildResourceCondition {
	return ChildResourceCondition{
		Created:      false,
		LastModified: time.Now().UTC().String(),
		Message:      "unable to proceed with resource creation " + err.Error(),
	}
}
