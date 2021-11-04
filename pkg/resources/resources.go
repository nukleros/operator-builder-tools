/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"
	"reflect"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/imdario/mergo"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToUnstructured returns an unstructured representation of a resource.
func ToUnstructured(resource metav1.Object) (*unstructured.Unstructured, error) {
	innerObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: innerObject}, nil
}

// ToProper returns the proper object representation of a resource.
func ToProper(destination metav1.Object, source metav1.Object) error {
	// ensure we are working with the same types
	if reflect.TypeOf(source) != reflect.TypeOf(destination) {
		return fmt.Errorf("type mismatch when converting to proper object")
	}

	// convert the source object to an unstructured type
	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(source)
	if err != nil {
		return err
	}

	// return the outcome of converting the unstructured type to its proper type
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObject, destination)
}

func getResourceChecker(resource metav1.Object) (resourceChecker, error) {
	runtimeObj, ok := resource.(runtime.Object)
	if !ok {
		return nil, fmt.Errorf("unable to convert metav1.Obect to runtime.Object")
	}

	switch runtimeObj.GetObjectKind().GroupVersionKind().Kind {
	case NamespaceKind:
		return NewNamespaceResource(resource)
	case CustomResourceDefinitionKind:
		return NewCRDResource(resource)
	case SecretKind:
		return NewSecretResource(resource)
	case ConfigMapKind:
		return NewConfigMapResource(resource)
	case DeploymentKind:
		return NewDeploymentResource(resource)
	case DaemonSetKind:
		return NewDaemonSetResource(resource)
	case StatefulSetKind:
		return NewStatefulSetResource(resource)
	case JobKind:
		return NewJobResource(resource)
	case ServiceKind:
		return NewServiceResource(resource)
	default:
		return NewUnknownResource(resource)
	}
}

// IsReady returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.
func IsReady(resource metav1.Object) (bool, error) {
	resourceChecker, err := getResourceChecker(resource)
	if err != nil {
		return false, fmt.Errorf("unable to determine ready status for resource, %w", err)

	}

	return resourceChecker.IsReady()
}

// AreReady returns whether resources are ready.  All resources must be ready in order
// to satisfy the requirement that resources are ready.
func AreReady(resources ...metav1.Object) (bool, error) {
	for _, rsrc := range resources {
		ready, err := IsReady(rsrc)
		if !ready || err != nil {
			return false, err
		}
	}

	return true, nil
}

// AreEqual determines if two resources are equal.
func AreEqual(desired, actual metav1.Object) (bool, error) {
	mergedResource, err := ToUnstructured(actual)
	if err != nil {
		return false, err
	}

	actualResource, err := ToUnstructured(actual)
	if err != nil {
		return false, err
	}

	desiredResource, err := ToUnstructured(desired)
	if err != nil {
		return false, err
	}

	// ensure that resource versions and observed generation do not interfere
	// with calculating equality
	desiredResource.SetResourceVersion(actualResource.GetResourceVersion())
	desiredResource.SetGeneration(actualResource.GetGeneration())

	// ensure that a current cluster-scoped resource is not evaluated against
	// a manifest which may include a namespace
	if actualResource.GetNamespace() == "" {
		desiredResource.SetNamespace(actualResource.GetNamespace())
	}

	// merge the overrides from the desired resource into the actual resource
	mergo.Merge(
		&mergedResource.Object,
		desiredResource.Object,
		mergo.WithOverride,
		mergo.WithSliceDeepCopy,
	)

	// calculate the actual differences
	diffOptions := []patch.CalculateOption{
		reconciler.IgnoreManagedFields(),
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
		patch.IgnorePDBSelector(),
	}

	diffResults, err := patch.DefaultPatchMaker.Calculate(
		actualResource,
		mergedResource,
		diffOptions...,
	)
	if err != nil {
		return false, err
	}

	return diffResults.IsEmpty(), nil
}

// EqualNamespaceName will compare the namespace and name of two resource objects for equality.
func EqualNamespaceName(left, right metav1.Object) bool {
	if left == nil || right == nil {
		return false
	}

	return (left.GetName() == right.GetName()) && (left.GetNamespace() == right.GetNamespace())
}

// EqualGVK will compare the GVK of two resource objects for equality.
func EqualGVK(left, right runtime.Object) bool {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false
	}

	return left.GetObjectKind().GroupVersionKind() == right.GetObjectKind().GroupVersionKind()
}
