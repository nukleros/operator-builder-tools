/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
)

const (
	ServiceKind = "Service"
)

// ServiceIsReady checks to see if a job is ready.
func ServiceIsReady(
	resource *Resource,
) (bool, error) {
	var service v1.Service
	if err := GetObject(resource, &service, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if service.Name == "" {
		return false, nil
	}

	// return if we have an external service type
	if service.Spec.Type == v1.ServiceTypeExternalName {
		return true, nil
	}

	// ensure a cluster ip address exists for cluster ip types
	if service.Spec.ClusterIP != v1.ClusterIPNone && service.Spec.ClusterIP == "" {
		return false, nil
	}

	// ensure a load balancer ip or hostname is present
	if service.Spec.Type == v1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) == 0 {
			return false, nil
		}
	}

	return true, nil
}
