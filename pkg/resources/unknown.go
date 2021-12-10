/*
	SPDX-License-Identifier: MIT
*/

package resources

import "sigs.k8s.io/controller-runtime/pkg/client"

// UnknownResource represents an unknown object.
type UnknownResource struct{}

// NewUnknownResource creates and returns a new UnknownResource.
func NewUnknownResource(object client.Object) (*UnknownResource, error) {
	return &UnknownResource{}, nil
}

// IsReady performs the logic to determine if an Unknown resource is ready.
func (unknown *UnknownResource) IsReady() (bool, error) {
	return true, nil
}
