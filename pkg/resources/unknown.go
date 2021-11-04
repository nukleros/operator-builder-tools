/*
	SPDX-License-Identifier: MIT
*/

package resources

type UnknownResource struct{}

// NewUnknownResource creates and returns a new UnknownResource.
func NewUnknownResource(name, namespace string) *UnknownResource {
	return &UnknownResource{}
}

// IsReady performs the logic to determine if an Unknown resource is ready.
func (unknown *UnknownResource) IsReady() (bool, error) {
	return true, nil
}
