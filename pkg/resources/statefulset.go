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
	StatefulSetKind    = "StatefulSet"
	StatefulSetVersion = "apps/v1"
)

type StatefulSetResource struct {
	parent *appsv1.StatefulSet
}

// NewStatefulSetResource creates and returns a new StatefulSetResource.
func NewStatefulSetResource(name, namespace string) *StatefulSetResource {
	return &StatefulSetResource{
		parent: &appsv1.StatefulSet{
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

// GetParent returns the parent attribute of the resource.
func (statefulSet *StatefulSetResource) GetParent() client.Object {
	return statefulSet.parent
}

// IsReady performs the logic to determine if a secret is ready.
func (statefulSet *StatefulSetResource) IsReady(resource *Resource) (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if statefulSet.parent.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if statefulSet.parent.Generation != statefulSet.parent.Status.ObservedGeneration {
		return false, nil
	}

	// check for valid replicas
	replicas := statefulSet.parent.Spec.Replicas
	if replicas == nil {
		return false, nil
	}

	// check to see if replicas have been updated
	var needsUpdate int32
	if statefulSet.parent.Spec.UpdateStrategy.RollingUpdate != nil &&
		statefulSet.parent.Spec.UpdateStrategy.RollingUpdate.Partition != nil &&
		*statefulSet.parent.Spec.UpdateStrategy.RollingUpdate.Partition > 0 {
		needsUpdate -= *statefulSet.parent.Spec.UpdateStrategy.RollingUpdate.Partition
	}

	notUpdated := needsUpdate - statefulSet.parent.Status.UpdatedReplicas
	if notUpdated > 0 {
		return false, nil
	}

	// check to see if replicas are available
	notReady := *replicas - statefulSet.parent.Status.ReadyReplicas
	if notReady > 0 {
		return false, nil
	}

	// check to see if a scale down operation is complete
	notDeleted := statefulSet.parent.Status.Replicas - *replicas
	if notDeleted > 0 {
		return false, nil
	}

	return true, nil
}
