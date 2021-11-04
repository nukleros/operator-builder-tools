/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CustomResourceDefinitionKind    = "CustomResourceDefinition"
	CustomResourceDefinitionVersion = "CustomResourceDefinition"
)

type CRDResource struct {
	extensionsv1.CustomResourceDefinition
}

// NewCRDResource creates and returns a new CRDResource.
func NewCRDResource(name, namespace string) *CRDResource {
	return &CRDResource{
		extensionsv1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				Kind:       CustomResourceDefinitionKind,
				APIVersion: CustomResourceDefinitionVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (crd *CRDResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if crd.Name == "" {
		return false, nil
	}

	return true, nil
}
