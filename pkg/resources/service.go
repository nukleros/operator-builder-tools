/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ServiceKind    = "Service"
	ServiceVersion = "v1"
)

// ServiceResource represents a Kubernetes Service object.
type ServiceResource struct {
	Object v1.Service
}

// NewServiceResource creates and returns a new ServiceResource.
func NewServiceResource(object client.Object) (*ServiceResource, error) {
	service := &v1.Service{}

	err := ToTyped(service, object)
	if err != nil {
		return nil, err
	}

	return &ServiceResource{Object: *service}, nil
}

// IsReady checks to see if a Service is ready.
func (service *ServiceResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if service.Object.Name == "" {
		return false, nil
	}

	// return if we have an external service type
	if service.Object.Spec.Type == v1.ServiceTypeExternalName {
		return true, nil
	}

	// ensure a cluster ip address exists for cluster ip types
	if service.Object.Spec.Type == v1.ServiceTypeClusterIP {
		if service.Object.Spec.ClusterIP == "" && len(service.Object.Spec.ClusterIPs) == 0 {
			return false, nil
		}
	}

	// ensure a load balancer ip or hostname is present
	if service.Object.Spec.Type == v1.ServiceTypeLoadBalancer {
		if len(service.Object.Status.LoadBalancer.Ingress) == 0 {
			return false, nil
		}
	}

	return true, nil
}
