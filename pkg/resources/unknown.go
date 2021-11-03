/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type UnknownResource struct{}

// NewUnknownResource creates and returns a new UnknownResource.
func NewUnknownResource() *UnknownResource {
	return &UnknownResource{}
}

// GetParent returns the parent attribute of the resource.
func (unknown *UnknownResource) GetParent() client.Object {
	return nil
}

// IsReady performs the logic to determine if a Unknown is ready.
func (unknown *UnknownResource) IsReady(resource *Resource) (bool, error) {
	return true, nil
}
