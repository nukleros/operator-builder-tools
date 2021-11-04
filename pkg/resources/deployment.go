/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DeploymentKind    = "Deployment"
	DeploymentVersion = "apps/v1"
)

type DeploymentResource struct {
	appsv1.Deployment
}

// NewDeploymentResource creates and returns a new DeploymentResource.
func NewDeploymentResource(name, namespace string) *DeploymentResource {
	return &DeploymentResource{
		appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       DeploymentKind,
				APIVersion: DeploymentVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// IsReady performs the logic to determine if a deployment is ready.
func (deployment *DeploymentResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if deployment.Name == "" {
		return false, nil
	}

	// if the object is equal to an empty object, we know we did not find the object
	if reflect.DeepEqual(deployment, &appsv1.Deployment{}) {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if deployment.Generation != deployment.Status.ObservedGeneration {
		return false, nil
	}

	// check the status for a ready deployment
	if deployment.Status.ReadyReplicas != deployment.Status.Replicas {
		return false, nil
	}

	return true, nil
}
