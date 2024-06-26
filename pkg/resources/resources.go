/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"
	"reflect"

	"github.com/cisco-open/k8s-objectmatcher/patch"
	"github.com/cisco-open/operator-tools/pkg/reconciler"

	"github.com/imdario/mergo"

	"github.com/nukleros/desired"
	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ToUnstructured returns an unstructured representation of a resource.
func ToUnstructured(resource client.Object) (*unstructured.Unstructured, error) {
	innerObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&resource)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: innerObject}, nil
}

// ToTyped returns the proper object representation of a resource.
func ToTyped(destination, source client.Object) error {
	// convert the source object to an unstructured type
	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(source)
	if err != nil {
		return err
	}

	// return the outcome of converting the unstructured type to its proper type
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObject, destination)
}

// getResourceCheckerFromReconciler gets a resource checker from an object.  This performs the resource
// lookup within the cluster from a reconciler.
func getResourceCheckerFromReconciler(r workload.Reconciler, req *workload.Request, resource client.Object) (resourceChecker, error) {
	if resource == nil {
		return nil, fmt.Errorf("no object was found")
	}

	clusterResource, err := Get(r, req, resource)
	if err != nil {
		return nil, err
	}

	switch resource.GetObjectKind().GroupVersionKind().Kind {
	case MutatingWebhookConfigurationKind:
		return NewMutatingWebhookConfigurationResource(r, req, clusterResource)
	case ValidatingWebhookConfigurationKind:
		return NewValidatingWebhookConfigurationResource(r, req, clusterResource)
	default:
		return getResourceChecker(clusterResource)
	}
}

// getResourceChecker gets a resource checker from an object.  This is only safe to assume that the
// object being passed has already been retrieved from the cluster.
func getResourceChecker(resource client.Object) (resourceChecker, error) {
	if resource == nil {
		return nil, fmt.Errorf("no object was found")
	}

	switch resource.GetObjectKind().GroupVersionKind().Kind {
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
	case EndpointsKind:
		return NewEndpointsResource(resource)
	case IssuerKind:
		return NewIssuerResource(resource)
	case ClusterIssuerKind:
		return NewClusterIssuerResource(resource)
	case CertificateKind:
		return NewCertificateResource(resource)
	default:
		return NewUnknownResource(resource)
	}
}

// IsReadyFromReconciler returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.  It takes in a client object and does the get on behalf of the caller.
func IsReadyFromReconciler(r workload.Reconciler, req *workload.Request, resource client.Object) (bool, error) {
	checker, err := getResourceCheckerFromReconciler(r, req, resource)
	if err != nil {
		return false, fmt.Errorf("unable to determine ready status for resource, %w", err)
	}

	return checker.IsReady()
}

// IsReady returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.
func IsReady(resource client.Object) (bool, error) {
	checker, err := getResourceChecker(resource)
	if err != nil {
		return false, fmt.Errorf("unable to determine ready status for resource, %w", err)
	}

	return checker.IsReady()
}

// AreReady returns whether resources are ready.  All resources must be ready in order
// to satisfy the requirement that resources are ready.
func AreReady(resources ...client.Object) (bool, error) {
	for i := range resources {
		ready, err := IsReady(resources[i])
		if !ready || err != nil {
			return false, err
		}
	}

	return true, nil
}

// AreEqual determines if two resources are equal.
func AreEqual(desired, actual client.Object) (bool, error) {
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
	if err := mergo.Merge(
		&mergedResource.Object,
		desiredResource.Object,
		mergo.WithOverride,
		mergo.WithSliceDeepCopy,
	); err != nil {
		return false, err
	}

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

// AreDesired determines if an actual resource is in a desired state based on the state
// of a desired resource.
func AreDesired(desiredObject, actualObject client.Object) (bool, error) {
	desiredResource, err := ToUnstructured(desiredObject)
	if err != nil {
		return false, err
	}

	actualResource, err := ToUnstructured(actualObject)
	if err != nil {
		return false, err
	}

	return desired.Desired(desiredResource.Object, actualResource.Object)
}

// EqualNamespaceName will compare the namespace and name of two resource objects for equality.
func EqualNamespaceName(left, right client.Object) bool {
	if left == nil || right == nil {
		return false
	}

	return (left.GetName() == right.GetName()) && (left.GetNamespace() == right.GetNamespace())
}

// EqualGVK will compare the GVK of two resource objects for equality.
func EqualGVK(left, right client.Object) bool {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return false
	}

	return left.GetObjectKind().GroupVersionKind() == right.GetObjectKind().GroupVersionKind()
}
