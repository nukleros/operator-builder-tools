/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DeploymentKind = "Deployment"
)

type DeploymentResource struct {
	parent *appsv1.Deployment
}

// NewDeploymentResource creates and returns a new DeploymentResource.
func NewDeploymentResource() *DeploymentResource {
	return &DeploymentResource{
		parent: &appsv1.Deployment{},
	}
}

// GetParent returns the parent attribute of the resource.
func (deployment *DeploymentResource) GetParent() client.Object {
	return deployment.parent
}

// IsReady performs the logic to determine if a deployment is ready.
func (deployment *DeploymentResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if deployment.parent.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if deployment.parent.Generation != deployment.parent.Status.ObservedGeneration {
		return false, nil
	}

	// check the status for a ready deployment
	if deployment.parent.Status.ReadyReplicas != deployment.parent.Status.Replicas {
		return false, nil
	}

	return true, nil
}
