/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DeploymentKind    = "Deployment"
	DeploymentVersion = "apps/v1"
)

type DeploymentResource struct {
	Object appsv1.Deployment
}

// NewDeploymentResource creates and returns a new DeploymentResource.
func NewDeploymentResource(object client.Object) (*DeploymentResource, error) {
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

	// check the status for a ready deployment
	if deployment.Object.Status.ReadyReplicas != deployment.Object.Status.Replicas {
		return false, nil
	}

	// ensure that there are no replicas that are unavailable
	if deployment.Object.Status.UnavailableReplicas > 0 {
		return false, nil
	}

	return true, nil
}
