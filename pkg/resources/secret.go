/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SecretKind    = "Secret"
	SecretVersion = "v1"
)

// SecretResource represents a Kubernetes Secret object.
type SecretResource struct {
	Object v1.Secret
}

// NewSecretResource creates and returns a new SecretResource.
func NewSecretResource(object client.Object) (*SecretResource, error) {
	secret := &v1.Secret{}

	err := ToTyped(secret, object)
	if err != nil {
		return nil, err
	}

	return &SecretResource{Object: *secret}, nil
}

// IsReady checks to see if a Secret is ready.
func (secret *SecretResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if secret.Object.Name == "" {
		return false, nil
	}

	return true, nil
}
