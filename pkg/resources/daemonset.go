/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DaemonSetKind    = "DaemonSet"
	DaemonSetVersion = "apps/v1"
)

type DaemonSetResource struct {
	Object appsv1.DaemonSet
}

// NewDaemonSetResource creates and returns a new DaemonSetResource.
func NewDaemonSetResource(object metav1.Object) (*DaemonSetResource, error) {
	daemonSet := &appsv1.DaemonSet{}

	err := ToProper(daemonSet, object)
	if err != nil {
		return nil, err
	}

	return &DaemonSetResource{Object: *daemonSet}, nil
}

// DaemonSetIsReady checks to see if a daemonset is ready.
func (daemonSet *DaemonSetResource) IsReady() (bool, error) {
	// ensure the desired number is scheduled and ready
	if daemonSet.Object.Status.DesiredNumberScheduled == daemonSet.Object.Status.NumberReady {
		if daemonSet.Object.Status.NumberReady > 0 && daemonSet.Object.Status.NumberUnavailable < 1 {
			return true, nil
		}
	}

	return false, nil
}
