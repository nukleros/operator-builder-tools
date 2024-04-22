/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ConfigMapKind    = "ConfigMap"
	ConfigMapVersion = "v1"
)

// ConfigMapResource represents a Kubernetes ConfigMap object.
type ConfigMapResource struct {
	Object v1.ConfigMap
}

// NewConfigMapResource creates and returns a new ConfigMapResource.
func NewConfigMapResource(object client.Object) (*ConfigMapResource, error) {
	configMap := &v1.ConfigMap{}

	err := ToTyped(configMap, object)
	if err != nil {
		return nil, err
	}

	return &ConfigMapResource{Object: *configMap}, nil
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (configMap *ConfigMapResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if configMap.Object.Name == "" {
		return false, nil
	}

	return true, nil
}
