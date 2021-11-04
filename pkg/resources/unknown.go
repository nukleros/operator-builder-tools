/*
	SPDX-License-Identifier: MIT
*/

package resources

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type UnknownResource struct{}

// NewUnknownResource creates and returns a new UnknownResource.
func NewUnknownResource(object metav1.Object) (*UnknownResource, error) {
	return &UnknownResource{}, nil
}

// IsReady performs the logic to determine if an Unknown resource is ready.
func (unknown *UnknownResource) IsReady() (bool, error) {
	return true, nil
}
