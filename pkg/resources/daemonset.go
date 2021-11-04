/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DaemonSetKind    = "DaemonSet"
	DaemonSetVersion = "apps/v1"
)

type DaemonSetResource struct {
	appsv1.DaemonSet
}

// NewDaemonSetResource creates and returns a new DaemonSetResource.
func NewDaemonSetResource(name, namespace string) *DaemonSetResource {
	return &DaemonSetResource{
		appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       DaemonSetKind,
				APIVersion: DaemonSetVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// GetParent returns the parent attribute of the resource.
func (daemonSet *DaemonSetResource) GetParent() client.Object {
	return daemonSet
}

// DaemonSetIsReady checks to see if a daemonset is ready.
func (daemonSet *DaemonSetResource) IsReady() (bool, error) {
	// ensure the desired number is scheduled and ready
	if daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady {
		if daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable < 1 {
			return true, nil
		}
	}

	return false, nil
}
