/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CustomResourceDefinitionKind    = "CustomResourceDefinition"
	CustomResourceDefinitionVersion = "CustomResourceDefinition"
)

type CRDResource struct {
	Object extensionsv1.CustomResourceDefinition
}

// NewCRDResource creates and returns a new CRDResource.
func NewCRDResource(object client.Object) (*CRDResource, error) {
	crd := &extensionsv1.CustomResourceDefinition{}

	if err := ToProper(crd, object); err != nil {
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
