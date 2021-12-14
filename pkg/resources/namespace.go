/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

const (
	NamespaceKind    = "Namespace"
	NamespaceVersion = "v1"
)

// NamespaceResource represents a Kubernetes Namespace object.
type NamespaceResource struct {
	Object v1.Namespace
}

// NewNamespaceResource creates and returns a new NamespaceResource.
func NewNamespaceResource(object client.Object) (*NamespaceResource, error) {
	namespace := &v1.Namespace{}

	err := ToTyped(namespace, object)
	if err != nil {
		return nil, err
	}

	return &NamespaceResource{Object: *namespace}, nil
}

// IsReady defines the criteria for a Namespace to be condsidered
// ready.
func (namespace *NamespaceResource) IsReady() (bool, error) {
	// if we have a name that is empty, we know we did not find the object
	if namespace.Object.Name == "" {
		return false, nil
	}

	// if the namespace is terminating, it is not considered ready
	if namespace.Object.Status.Phase == v1.NamespaceTerminating {
		return false, nil
	}

	// finally, rely on the active field to determine if this namespace is ready
	return namespace.Object.Status.Phase == v1.NamespaceActive, nil
}

// NamespaceForResourceIsReady will check to see if the Namespace of a metadata.namespace
// field of a resource is ready.
func NamespaceForResourceIsReady(r workload.Reconciler, req *workload.Request, resource client.Object) (bool, error) {
	namespace := &v1.Namespace{}
	namespacedName := types.NamespacedName{
		Name:      resource.GetNamespace(),
		Namespace: "",
	}

	if err := r.Get(req.Context, namespacedName, namespace); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("unable to get namespace, %w", err)
	}

	return IsReady(namespace)
}
