/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StatefulSetKind    = "StatefulSet"
	StatefulSetVersion = "apps/v1"
)

type StatefulSetResource struct {
	appsv1.StatefulSet
}

// NewStatefulSetResource creates and returns a new StatefulSetResource.
func NewStatefulSetResource(name, namespace string) *StatefulSetResource {
	return &StatefulSetResource{
		appsv1.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       StatefulSetKind,
				APIVersion: StatefulSetVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

// IsReady performs the logic to determine if a secret is ready.
func (statefulSet *StatefulSetResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if statefulSet.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if statefulSet.Generation != statefulSet.Status.ObservedGeneration {
		return false, nil
	}

	// check for valid replicas
	replicas := statefulSet.Spec.Replicas
	if replicas == nil {
		return false, nil
	}

	// check to see if replicas have been updated
	var needsUpdate int32
	if statefulSet.Spec.UpdateStrategy.RollingUpdate != nil &&
		statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition != nil &&
		*statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition > 0 {
		needsUpdate -= *statefulSet.Spec.UpdateStrategy.RollingUpdate.Partition
	}

	notUpdated := needsUpdate - statefulSet.Status.UpdatedReplicas
	if notUpdated > 0 {
		return false, nil
	}

	// check to see if replicas are available
	notReady := *replicas - statefulSet.Status.ReadyReplicas
	if notReady > 0 {
		return false, nil
	}

	// check to see if a scale down operation is complete
	notDeleted := statefulSet.Status.Replicas - *replicas
	if notDeleted > 0 {
		return false, nil
	}

	return true, nil
}
