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
	Object extensionsv1.CustomResourceDefinition
}

// NewCRDResource creates and returns a new CRDResource.
func NewCRDResource(object metav1.Object) (*CRDResource, error) {
	crd := &extensionsv1.CustomResourceDefinition{}

	err := ToProper(crd, object)
	if err != nil {
		return nil, err
	}

	return &CRDResource{Object: *crd}, nil
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (crd *CRDResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if crd.Object.Name == "" {
		return false, nil
	}

	return true, nil
}
