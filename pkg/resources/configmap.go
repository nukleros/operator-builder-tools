/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ConfigMapKind = "ConfigMap"
)

type ConfigMapResource struct {
	parent *v1.ConfigMap
}

// NewConfigMapResource creates and returns a new ConfigMapResource.
func NewConfigMapResource() *ConfigMapResource {
	return &ConfigMapResource{
		parent: &v1.ConfigMap{},
	}
}

// GetParent returns the parent attribute of the resource.
func (configMap *ConfigMapResource) GetParent() client.Object {
	return configMap.parent
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (configMap *ConfigMapResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if configMap.parent.Name == "" {
		return false, nil
	}

	return true, nil
}
