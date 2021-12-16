// SPDX-License-Identifier: MIT

package workload

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/status"
)

var ErrCollectionNotFound = errors.New("collection not found")
var ErrInvalidWorkload = errors.New("supplied workload is invalid")

// Workload represents a Custom Resource that your controller is watching.
type Workload interface {
	client.Object

	GetWorkloadGVK() schema.GroupVersionKind
	GetDependencies() []Workload
	GetDependencyStatus() bool
	GetReadyStatus() bool
	GetPhaseConditions() []*status.PhaseCondition
	GetChildResourceConditions() []*status.ChildResource

	SetReadyStatus(bool)
	SetDependencyStatus(bool)
	SetPhaseCondition(*status.PhaseCondition)
	SetChildResourceCondition(*status.ChildResource)
}

// Validate validates an individual workload to ensure that its GVK is for the
// correct resource.
func Validate(workload Workload) error {
	defaultWorkloadGVK := workload.GetWorkloadGVK()

	if defaultWorkloadGVK != workload.GetObjectKind().GroupVersionKind() {
		return fmt.Errorf(
			"%w, expected resource of kind: '%s', with group '%s' and version '%s'; "+
				"found resource of kind '%s', with group '%s' and version '%s'",
			ErrInvalidWorkload,
			defaultWorkloadGVK.Kind,
			defaultWorkloadGVK.Group,
			defaultWorkloadGVK.Version,
			workload.GetObjectKind().GroupVersionKind().Kind,
			workload.GetObjectKind().GroupVersionKind().Group,
			workload.GetObjectKind().GroupVersionKind().Version,
		)
	}

	return nil
}
