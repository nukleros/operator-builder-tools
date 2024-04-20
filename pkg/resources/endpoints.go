/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	EndpointsKind    = "Endpoints"
	EndpointsVersion = "v1"
)

type EndpointsResource struct {
	Object corev1.Endpoints
}

// NewEndpointsResource creates and returns a new EndpointResource.
func NewEndpointsResource(object client.Object) (*EndpointsResource, error) {
	endpoints := &corev1.Endpoints{}

	err := ToTyped(endpoints, object)
	if err != nil {
		return nil, err
	}

	return &EndpointsResource{Object: *endpoints}, nil
}

// IsReady performs the logic to determine if an Endpoints resource is ready.
func (endpoints *EndpointsResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if endpoints.Object.Name == "" {
		return false, nil
	}

	return len(endpoints.Object.Subsets) > 0, nil
}
