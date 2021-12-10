// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

// LifecycleEvent is used to convey which lifecycle event we are targeting.
type LifecycleEvent int32

const (
	CreateEvent LifecycleEvent = iota
	UpdateEvent
	DeleteEvent
)

// Registry is a store for all the phases for each event loop.
type Registry struct {
	createPhases []*Phase
	updatePhases []*Phase
	deletePhases []*Phase
}

// Register is used to add a Phase to the Registry for the provided event loop.
func (registry *Registry) Register(name string, definition HandlerFunc, event LifecycleEvent, options ...PhaseOption) {
	phase := &Phase{
		Name:       name,
		definition: definition,
	}

	for _, option := range options {
		option(phase)
	}

	switch event {
	case CreateEvent:
		registry.createPhases = append(registry.createPhases, phase)
	case UpdateEvent:
		registry.updatePhases = append(registry.updatePhases, phase)
	case DeleteEvent:
		registry.deletePhases = append(registry.deletePhases, phase)
	}
}

// HandleExecution will trigger the execution of the phases
// for the appropriate lifecycle event. This is the main entrypoint
// into our phases.
func (registry *Registry) HandleExecution(r workload.Reconciler, req *workload.Request) (reconcile.Result, error) {
	// execute the phases
	switch {
	case !req.Workload.GetDeletionTimestamp().IsZero():
		req.Log.Info("deleting workload")
		myFinalizerName := fmt.Sprintf("%s/Finalizer", req.Workload.GetWorkloadGVK().Group)

		// The object is being deleted
		if containsString(req.Workload.GetFinalizers(), myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			result, err := registry.Execute(r, req, DeleteEvent)
			if err != nil || !result.IsZero() {
				return result, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(req.Workload, myFinalizerName)

			if err := r.Update(req.Context, req.Workload); err != nil {
				return ctrl.Result{}, fmt.Errorf("unable to remove finalizer from %s, %w", req.Workload.GetWorkloadGVK().Kind, err)
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	case !req.Workload.GetReadyStatus():
		return registry.Execute(r, req, CreateEvent)
	default:
		return registry.Execute(r, req, UpdateEvent)
	}
}

// Execute runs the phases for the specified lifecycle event.
func (registry *Registry) Execute(r workload.Reconciler, req *workload.Request, event LifecycleEvent) (reconcile.Result, error) {
	phases := registry.getPhases(event)
	for _, phase := range phases {
		req.Log.V(5).Info(
			"enter phase",
			"phase", phase.Name,
		)

		proceed, err := phase.definition(r, req)
		result, err := phase.handlePhaseExit(r, req, proceed, err)

		if err != nil || !proceed {
			req.Log.V(2).Info(
				"not ready; requeuing",
				"phase", phase.Name,
			)

			// return only if we have an error or are told not to proceed
			if err != nil {
				return result, fmt.Errorf("unable to complete %s phase for %s, %w", phase.Name, req.Workload.GetWorkloadGVK().Kind, err)
			}

			if !proceed {
				return result, nil
			}
		}

		req.Log.V(5).Info(
			"completed phase",
			"phase", phase.Name,
		)
	}

	return ctrl.Result{}, nil
}

// getPhases returns the phases for a given lifecycle event.
func (registry *Registry) getPhases(event LifecycleEvent) []*Phase {
	switch event {
	case CreateEvent:
		return registry.createPhases
	case UpdateEvent:
		return registry.updatePhases
	case DeleteEvent:
		return registry.deletePhases
	}

	return nil
}
