/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SecretKind    = "Secret"
	SecretVersion = "v1"
)

type SecretResource struct {
	parent *v1.Secret
}

// NewSecretResource creates and returns a new SecretResource.
func NewSecretResource(name, namespace string) *SecretResource {
	return &SecretResource{
		parent: &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       SecretKind,
				APIVersion: SecretVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// GetParent returns the parent attribute of the resource.
func (secret *SecretResource) GetParent() client.Object {
	return secret.parent
}

// IsReady checks to see if a secret is ready.
func (secret *SecretResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if secret.parent.Name == "" {
		return false, nil
	}

	return true, nil
}
