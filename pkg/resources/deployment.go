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
	Object appsv1.Deployment
}

// NewDeploymentResource creates and returns a new DeploymentResource.
func NewDeploymentResource(object metav1.Object) (*DeploymentResource, error) {
	deployment := &appsv1.Deployment{}

	err := ToProper(deployment, object)
	if err != nil {
		return nil, err
	}

	return &DeploymentResource{Object: *deployment}, nil
}

// IsReady performs the logic to determine if a deployment is ready.
func (deployment *DeploymentResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if deployment.Object.Name == "" {
		return false, nil
	}

	// if the object is equal to an empty object, we know we did not find the object
	if reflect.DeepEqual(deployment, &appsv1.Deployment{}) {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if deployment.Object.Generation != deployment.Object.Status.ObservedGeneration {
		return false, nil
	}

	// check the status for a ready deployment
	if deployment.Object.Status.ReadyReplicas != deployment.Object.Status.Replicas {
		return false, nil
	}

	return true, nil
}
