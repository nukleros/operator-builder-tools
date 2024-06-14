/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ReadyPathAnnotation  = "operator-builder.nukleros.io/ready-path"
	ReadyValueAnnotation = "operator-builder.nukleros.io/ready-value"
)

// UnknownResource represents an unknown object.
type UnknownResource struct {
	Object client.Object
}

// NewUnknownResource creates and returns a new UnknownResource.
func NewUnknownResource(object client.Object) (*UnknownResource, error) {
	return &UnknownResource{Object: object}, nil
}

// IsReady performs the logic to determine if an Unknown resource is ready.  It allows functionality
// that will read a set of annotations that take in a field (in JSONpath format) and a value.  If
// the specific field is equal to that value, then the resource is considered to be ready.
func (unknown *UnknownResource) IsReady() (bool, error) {
	return isReadyFromAnnotations(unknown.Object)
}

// isReadyFromAnnotations checks to see if the specific resource has annotations indicating that
// the consumer is looking for a specific value from a specific path to determine resource
// readiness.  It is intended to mostly default to true unless a user messed up their
// JSONpath input.
func isReadyFromAnnotations(object client.Object) (bool, error) {
	if object == nil {
		return true, nil
	}

	annotations := object.GetAnnotations()
	if annotations == nil {
		return true, nil
	}

	// get inputs
	path := annotations[ReadyPathAnnotation]
	value := annotations[ReadyValueAnnotation]

	if path == "" || value == "" {
		return true, nil
	}

	asUnstructured, err := ToUnstructured(object)
	if err != nil {
		return false, err
	}

	// check if the JSONPath query matches the value and return
	ready, err := matchJSONPath(asUnstructured, path, value)
	if err != nil {
		return false, fmt.Errorf(
			"unable to determine resource readiness for path annotation [%s] and value annotation [%s] - %w",
			ReadyPathAnnotation,
			ReadyValueAnnotation,
			err,
		)
	}

	return ready, nil
}

// matchJSONPath checks if the JSONPath query `path` matches the given `value` in the Kubernetes object `obj`.
func matchJSONPath(object *unstructured.Unstructured, path, value string) (bool, error) {
	query := jsonpath.New("jsonpath")

	if err := query.Parse(fmt.Sprintf("{%s}", path)); err != nil {
		return false, fmt.Errorf("unable to parse query string [%s] - %w", path, err)
	}

	results, err := query.FindResults(object.Object)
	if err != nil {
		// if part of the path is not found, we simply return that the resource is not ready
		if strings.Contains(err.Error(), "is not found") {
			return false, nil
		}

		return false, fmt.Errorf("unable to find results for query string [%s] and value [%s] - %w", path, value, err)
	}

	if len(results) == 0 || len(results[0]) == 0 {
		return false, nil
	}

	return strings.TrimSpace(fmt.Sprintf("%v", results[0][0].Interface())) == value, nil
}
