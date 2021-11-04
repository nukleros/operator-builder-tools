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
	ConfigMapKind    = "ConfigMap"
	ConfigMapVersion = "v1"
)

type ConfigMapResource struct {
	v1.ConfigMap
}

// NewConfigMapResource creates and returns a new ConfigMapResource.
func NewConfigMapResource(name, namespace string) *ConfigMapResource {
	return &ConfigMapResource{
		v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       ConfigMapKind,
				APIVersion: ConfigMapVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// GetParent returns the parent attribute of the resource.
func (configMap *ConfigMapResource) GetParent() client.Object {
	return configMap
}

// IsReady performs the logic to determine if a ConfigMap is ready.
func (configMap *ConfigMapResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if configMap.Name == "" {
		return false, nil
	}

	return true, nil
}
