// SPDX-License-Identifier: MIT

package workload

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Request holds the state of the current reconcile request
// This Object can be used to pass state such as context
// between phases of the controller.
type Request struct {
	Context    context.Context
	Workload   Workload
	Collection Workload
	Resources  []client.Object
	Log        logr.Logger
}
