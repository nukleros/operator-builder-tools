/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
)

const (
	ConfigMapKind = "ConfigMap"
)

// ConfigMapIsReady performs the logic to determine if a secret is ready.
func ConfigMapIsReady(
	resource *Resource,
	expectedKeys ...string,
) (bool, error) {
	var configMap v1.ConfigMap
	if err := GetObject(resource, &configMap, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if configMap.Name == "" {
		return false, nil
	}

	// check that expected keys are set in the configmap
	for _, key := range expectedKeys {
		if configMap.Data[key] == "" {
			return false, nil
		}
	}

	return true, nil
}
