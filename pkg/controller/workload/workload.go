// SPDX-License-Identifier: MIT

package workload

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/status"
)

var ErrCollectionNotFound = errors.New("collection not found")

// Workload represents a Custom Resource that your controller is watching.
type Workload interface {
	client.Object
	GetWorkloadGVK() schema.GroupVersionKind
	GetDependencies() []Workload
	GetDependencyStatus() bool
	GetReadyStatus() bool
	GetPhaseConditions() []*status.PhaseCondition
	GetResourceConditions() []*status.ChildResource

	SetReadyStatus(bool)
	SetDependencyStatus(bool)
	SetPhaseCondition(*status.PhaseCondition)
	SetResourceCondition(*status.ChildResource)
}
