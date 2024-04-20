/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DaemonSetKind    = "DaemonSet"
	DaemonSetVersion = "apps/v1"
)

// DaemonSetResource represents a Kubernetes Daemonset object.
type DaemonSetResource struct {
	Object appsv1.DaemonSet
}

// NewDaemonSetResource creates and returns a new DaemonSetResource.
func NewDaemonSetResource(object client.Object) (*DaemonSetResource, error) {
	daemonSet := &appsv1.DaemonSet{}

	err := ToTyped(daemonSet, object)
	if err != nil {
		return nil, err
	}

	return &DaemonSetResource{Object: *daemonSet}, nil
}

// IsReady checks to see if a DaemonSet is ready.
func (daemonSet *DaemonSetResource) IsReady() (bool, error) {
	// ensure the desired number is scheduled and ready
	if daemonSet.Object.Status.DesiredNumberScheduled == daemonSet.Object.Status.NumberReady {
		if daemonSet.Object.Status.NumberReady > 0 && daemonSet.Object.Status.NumberUnavailable < 1 {
			return true, nil
		}
	}

	return false, nil
}
