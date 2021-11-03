/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	CustomResourceDefinitionKind = "CustomResourceDefinition"
)

type CRDResource struct {
	parent *extensionsv1.CustomResourceDefinition
}

// NewCRDResource creates and returns a new CRDResource.
func NewCRDResource() *CRDResource {
	return &CRDResource{
		parent: &extensionsv1.CustomResourceDefinition{},
	}
}

// GetParent returns the parent attribute of the resource.
func (crd *CRDResource) GetParent() client.Object {
	return crd.parent
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (crd *CRDResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if crd.parent.Name == "" {
		return false, nil
	}

	return true, nil
}
