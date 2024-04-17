// SPDX-License-Identifier: MIT

package workload

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Reconciler acts as a client and state holder of the controller
// for use within the rconciliation loop.
type Reconciler interface {
	client.Client

	// attribute exporters and setters
	GetController() controller.Controller
	GetManager() manager.Manager
	GetLogger() logr.Logger
	GetResources(*Request) ([]client.Object, error)
	GetEventRecorder() record.EventRecorder
	GetFieldManager() string
	GetWatches() []client.Object
	SetWatch(client.Object)

	// custom methods which are managed by consumers
	CheckReady(*Request) (bool, error)
	Mutate(*Request, client.Object) ([]client.Object, bool, error)
}
