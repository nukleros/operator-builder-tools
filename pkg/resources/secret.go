/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SecretKind = "Secret"
)

type SecretResource struct {
	parent *v1.Secret
}

// NewSecretResource creates and returns a new SecretResource.
func NewSecretResource() *SecretResource {
	return &SecretResource{
		parent: &v1.Secret{},
	}
}

// GetParent returns the parent attribute of the resource.
func (secret *SecretResource) GetParent() client.Object {
	return secret.parent
}

// IsReady checks to see if a secret is ready.
func (secret *SecretResource) IsReady(resource *resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if secret.parent.Name == "" {
		return false, nil
	}

	return true, nil
}
