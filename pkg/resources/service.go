/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ServiceKind = "Service"
)

type ServiceResource struct {
	parent *v1.Service
}

// NewServiceResource creates and returns a new ServiceResource.
func NewServiceResource() *ServiceResource {
	return &ServiceResource{
		parent: &v1.Service{},
	}
}

// GetParent returns the parent attribute of the resource.
func (service *ServiceResource) GetParent() client.Object {
	return service.parent
}

// IsReady checks to see if a job is ready.
func (service *ServiceResource) IsReady(resource *resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if service.parent.Name == "" {
		return false, nil
	}

	// return if we have an external service type
	if service.parent.Spec.Type == v1.ServiceTypeExternalName {
		return true, nil
	}

	// ensure a cluster ip address exists for cluster ip types
	if service.parent.Spec.ClusterIP != v1.ClusterIPNone && service.parent.Spec.ClusterIP == "" {
		return false, nil
	}

	// ensure a load balancer ip or hostname is present
	if service.parent.Spec.Type == v1.ServiceTypeLoadBalancer {
		if len(service.parent.Status.LoadBalancer.Ingress) == 0 {
			return false, nil
		}
	}

	return true, nil
}
