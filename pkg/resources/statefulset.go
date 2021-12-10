/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	StatefulSetKind    = "StatefulSet"
	StatefulSetVersion = "apps/v1"
)

// StatefulSetResource represents a Kubernetes StatefulSet object.
type StatefulSetResource struct {
	Object appsv1.StatefulSet
}

// NewStatefulSetResource creates and returns a new StatefulSetResource.
func NewStatefulSetResource(object client.Object) (*StatefulSetResource, error) {
	statefulSet := &appsv1.StatefulSet{}

	err := ToTyped(statefulSet, object)
	if err != nil {
		return nil, err
	}

	return &StatefulSetResource{Object: *statefulSet}, nil
}

// IsReady performs the logic to determine if a StatefulSet is ready.
func (statefulSet *StatefulSetResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if statefulSet.Object.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if statefulSet.Object.Generation != statefulSet.Object.Status.ObservedGeneration {
		return false, nil
	}

	// check for valid replicas
	replicas := statefulSet.Object.Spec.Replicas
	if replicas == nil {
		return false, nil
	}

	// check to see if replicas have been updated
	var needsUpdate int32
	if statefulSet.Object.Spec.UpdateStrategy.RollingUpdate != nil &&
		statefulSet.Object.Spec.UpdateStrategy.RollingUpdate.Partition != nil &&
		*statefulSet.Object.Spec.UpdateStrategy.RollingUpdate.Partition > 0 {
		needsUpdate -= *statefulSet.Object.Spec.UpdateStrategy.RollingUpdate.Partition
	}

	notUpdated := needsUpdate - statefulSet.Object.Status.UpdatedReplicas
	if notUpdated > 0 {
		return false, nil
	}

	// check to see if replicas are available
	notReady := *replicas - statefulSet.Object.Status.ReadyReplicas
	if notReady > 0 {
		return false, nil
	}

	// check to see if a scale down operation is complete
	notDeleted := statefulSet.Object.Status.Replicas - *replicas
	if notDeleted > 0 {
		return false, nil
	}

	return true, nil
}
