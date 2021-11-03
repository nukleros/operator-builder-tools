/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/imdario/mergo"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const (
	FieldManager = "reconciler"
)

// ToUnstructured returns an unstructured representation of a Resource.
func (resource *Resource) ToUnstructured() (*unstructured.Unstructured, error) {
	innerObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&resource.Object)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: innerObject}, nil
}

// ToCommonResource converts a resources.Resource into a common API resource.
func (resource *Resource) ToCommonResource() *Resource {
	commonResource := &Resource{}

	// set the inherited fields
	commonResource.Group = resource.Group
	commonResource.Version = resource.Version
	commonResource.Kind = resource.Kind
	commonResource.Name = resource.Name
	commonResource.Namespace = resource.Namespace

	return commonResource
}

func (resource *Resource) setResourceChecker() {
	switch resource.Kind {
	case NamespaceKind:
		resource.resourceChecker = NewNamespaceResource()
	case CustomResourceDefinitionKind:
		resource.resourceChecker = NewCRDResource()
	case SecretKind:
		resource.resourceChecker = NewSecretResource()
	case ConfigMapKind:
		resource.resourceChecker = NewConfigMapResource()
	case DeploymentKind:
		resource.resourceChecker = NewDeploymentResource()
	case DaemonSetKind:
		resource.resourceChecker = NewDaemonSetResource()
	case StatefulSetKind:
		resource.resourceChecker = NewStatefulSetResource()
	case JobKind:
		resource.resourceChecker = NewJobResource()
	case ServiceKind:
		resource.resourceChecker = NewServiceResource()
	default:
		resource.resourceChecker = NewUnknownResource()
	}
}

// IsReady returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.
func (resource *Resource) IsReady() (bool, error) {
	resource.setResourceChecker()

	// get the object from the kubernetes cluster
	if err := GetObject(resource, true); err != nil {
		return false, err
	}

	return resource.resourceChecker.IsReady(resource)
}

// AreReady returns whether resources are ready.  All resources must be ready in order
// to satisfy the requirement that resources are ready.
func AreReady(resources ...*Resource) (bool, error) {
	for _, resource := range resources {
		ready, err := resource.IsReady()
		if !ready || err != nil {
			return false, err
		}
	}

	return true, nil
}

// AreEqual determines if two resources are equal.
func AreEqual(desired, actual *Resource) (bool, error) {
	mergedResource, err := actual.ToUnstructured()
	if err != nil {
		return false, err
	}

	actualResource, err := actual.ToUnstructured()
	if err != nil {
		return false, err
	}

	desiredResource, err := desired.ToUnstructured()
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
func (resource *Resource) EqualNamespaceName(compared *Resource) bool {
	return (resource.Name == compared.Name) && (resource.Namespace == compared.Namespace)
}

// EqualGVK will compare the GVK of two resource objects for equality.
func (resource *Resource) EqualGVK(compared *Resource) bool {
	return resource.Group == compared.Group &&
		resource.Version == compared.Version &&
		resource.Kind == compared.Kind
}

// GetObject returns an object based on an input object, and a destination
// object.
func GetObject(
	source *Resource,
	allowMissing bool,
) error {
	namespacedName := types.NamespacedName{
		Name:      source.Name,
		Namespace: source.Namespace,
	}

	if err := source.Client.Get(
		source.Context,
		namespacedName,
		source.resourceChecker.GetParent(),
	); err != nil {
		if allowMissing {
			if errors.IsNotFound(err) {
				return nil
			}
		} else {
			return fmt.Errorf(
				"unable to fetch resource of kind: [%s] in namespaced name: [%v]; %w",
				source.Kind,
				namespacedName,
				err,
			)
		}
	}

	return nil
}
