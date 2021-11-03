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

// resource represents any kubernetes resource.
type Resource struct {
	ResourceCommon
	resourceChecker resourceChecker

	Client  client.Client
	Context context.Context
	Object  client.Object
}

// resourceChecker is an interface which allows checking of a resource to see
// if it is in a ready state.
type resourceChecker interface {
	IsReady(*Resource) (bool, error)
	GetParent() client.Object
}

// NewResource returns a new resource given a client object and a kubernetes api client
// to use for interacting with cluster objects.
func NewResource(object client.Object, apiClient client.Client, ctx context.Context) *Resource {
	newResource := &Resource{
		Object:  object,
		Client:  apiClient,
		Context: ctx,
	}

	// set the inherited fields
	newResource.Group = object.GetObjectKind().GroupVersionKind().Group
	newResource.Version = object.GetObjectKind().GroupVersionKind().Version
	newResource.Kind = object.GetObjectKind().GroupVersionKind().Kind
	newResource.Name = object.GetName()
	newResource.Namespace = object.GetNamespace()

	// set the resourceChecker object
	newResource.setResourceChecker()

	return newResource
}
