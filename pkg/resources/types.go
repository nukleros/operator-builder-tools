/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceCommon are the common fields used across multiple resource types.
type ResourceCommon struct {
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
}

// Resource represents any kubernetes resource.
type Resource struct {
	ResourceCommon

	Object  client.Object
	Client  client.Client
	Context context.Context
	// Reconciler common.ComponentReconciler
}
