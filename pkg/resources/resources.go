/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/imdario/mergo"

	"sigs.k8s.io/controller-runtime/pkg/client"

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

// IsReady returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.
func (resource *Resource) IsReady() (bool, error) {
	switch resource.Kind {
	case NamespaceKind:
		return NamespaceIsReady(resource)
	case CustomResourceDefinitionKind:
		return CustomResourceDefinitionIsReady(resource)
	case SecretKind:
		return SecretIsReady(resource)
	case ConfigMapKind:
		return ConfigMapIsReady(resource)
	case DeploymentKind:
		return DeploymentIsReady(resource)
	case DaemonSetKind:
		return DaemonSetIsReady(resource)
	case StatefulSetKind:
		return StatefulSetIsReady(resource)
	case JobKind:
		return JobIsReady(resource)
	case ServiceKind:
		return ServiceIsReady(resource)
	}

	return true, nil
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

// NeedsUpdate determines if a resource needs to be updated.
func NeedsUpdate(desired, actual *Resource) (bool, error) {
	// check for equality first as this will let us avoid spamming user logs
	// when resources that need to be skipped explicitly (e.g. CRDs) are seen
	// as equal anyway
	equal, err := AreEqual(desired, actual)
	if equal || err != nil {
		return !equal, err
	}

	// // always skip custom resource updates as they are sensitive to modification
	// // e.g. resources provisioned by the resource definition would not
	// // understand the update to a spec
	// if desired.Kind == "CustomResourceDefinition" {
	// 	message := fmt.Sprintf("skipping update of CustomResourceDefinition "+
	// 		"[%s]", desired.Name)
	// 	messageVerbose := fmt.Sprintf("if updates to CustomResourceDefinition "+
	// 		"[%s] are desired, consider re-deploying the parent "+
	// 		"resource or generating a new api version with the desired "+
	// 		"changes", desired.Name)
	// 	desired.Reconciler.GetLogger().V(4).Info(message)
	// 	desired.Reconciler.GetLogger().V(7).Info(messageVerbose)

	// 	return false, nil
	// }

	return true, nil
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
	destination client.Object,
	allowMissing bool,
) error {
	namespacedName := types.NamespacedName{
		Name:      source.Name,
		Namespace: source.Namespace,
	}

	if err := source.Client.Get(
		source.Context,
		namespacedName,
		destination,
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
