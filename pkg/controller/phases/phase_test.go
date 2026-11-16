/*
	SPDX-License-Identifier: MIT
*/

package phases

import (
	"context"
	"errors"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	"github.com/nukleros/operator-builder-tools/pkg/status"
)

var testGVK = schema.GroupVersionKind{Group: "test.nukleros.io", Version: "v1", Kind: "Test"}

// testWorkload is a workload backed by an unstructured object, with its status kept in memory.
type testWorkload struct {
	*unstructured.Unstructured

	conditions []*status.PhaseCondition
}

func newTestWorkload(deleting bool) *testWorkload {
	w := &testWorkload{Unstructured: &unstructured.Unstructured{}}
	w.SetGroupVersionKind(testGVK)
	w.SetName("test")
	w.SetFinalizers([]string{testGVK.Group + "/Finalizer"})

	if deleting {
		now := metav1.Now()
		w.SetDeletionTimestamp(&now)
	}

	return w
}

func (w *testWorkload) GetWorkloadGVK() schema.GroupVersionKind             { return testGVK }
func (w *testWorkload) GetDependencies() []workload.Workload                { return nil }
func (w *testWorkload) GetDependencyStatus() bool                           { return true }
func (w *testWorkload) GetReadyStatus() bool                                { return true }
func (w *testWorkload) GetPhaseConditions() []*status.PhaseCondition        { return w.conditions }
func (w *testWorkload) GetChildResourceConditions() []*status.ChildResource { return nil }
func (w *testWorkload) SetReadyStatus(bool)                                 {}
func (w *testWorkload) SetDependencyStatus(bool)                            {}
func (w *testWorkload) SetChildResourceCondition(*status.ChildResource)     {}
func (w *testWorkload) SetPhaseCondition(condition *status.PhaseCondition) {
	w.conditions = append(w.conditions, condition)
}

// testReconciler is a reconciler whose updates of the workload, and of its status, return the
// configured errors.  Anything else it is asked to do panics on the nil embedded reconciler.
type testReconciler struct {
	workload.Reconciler

	updateErr       error
	statusUpdateErr error
}

func (r *testReconciler) Update(context.Context, client.Object, ...client.UpdateOption) error {
	return r.updateErr
}

func (r *testReconciler) Status() client.SubResourceWriter {
	return testStatusWriter{err: r.statusUpdateErr}
}

type testStatusWriter struct {
	client.SubResourceWriter

	err error
}

func (w testStatusWriter) Update(context.Context, client.Object, ...client.SubResourceUpdateOption) error {
	return w.err
}

func errNotFound() error {
	return apierrors.NewNotFound(schema.GroupResource{Group: testGVK.Group, Resource: "tests"}, "test")
}

func TestHandlePhaseExit(t *testing.T) {
	t.Parallel()

	errPhase := errors.New("phase failed")
	errUpdate := errors.New("update failed")

	tests := []struct {
		name            string
		deleting        bool
		phaseErr        error
		statusUpdateErr error
		wantErr         bool
		wantPhaseErr    bool
	}{
		{
			name: "a successful status update returns no error",
		},
		{
			name:            "a status update that fails returns its error",
			statusUpdateErr: errUpdate,
			wantErr:         true,
		},
		{
			name:            "a workload that is gone while being deleted has no status left to update",
			deleting:        true,
			statusUpdateErr: errNotFound(),
		},
		{
			name:            "a workload that is gone while being deleted still returns the phase error",
			deleting:        true,
			phaseErr:        errPhase,
			statusUpdateErr: errNotFound(),
			wantErr:         true,
			wantPhaseErr:    true,
		},
		{
			name:            "a workload that is gone while it is not being deleted is an error",
			statusUpdateErr: errNotFound(),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			phase := &Phase{Name: "Test"}
			req := &workload.Request{Context: context.Background(), Workload: newTestWorkload(tt.deleting)}
			r := &testReconciler{statusUpdateErr: tt.statusUpdateErr}

			_, err := phase.handlePhaseExit(r, req, tt.phaseErr == nil, tt.phaseErr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("handlePhaseExit() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantPhaseErr && !errors.Is(err, errPhase) {
				t.Errorf("handlePhaseExit() error = %v, want the phase error", err)
			}
		})
	}
}
