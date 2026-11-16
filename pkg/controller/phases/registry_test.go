/*
	SPDX-License-Identifier: MIT
*/

package phases

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
)

func TestHandleExecutionDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		updateErr       error
		statusUpdateErr error
		wantErr         bool
	}{
		{
			name: "the finalizer is removed",
		},
		{
			name:            "a workload another reconcile already deleted is not an error",
			updateErr:       errNotFound(),
			statusUpdateErr: errNotFound(),
		},
		{
			name:      "a finalizer that cannot be removed is an error",
			updateErr: errors.New("update failed"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registry := &Registry{}
			registry.Register("DeletionComplete", DeletionCompletePhase, DeleteEvent)

			req := &workload.Request{
				Context:  context.Background(),
				Workload: newTestWorkload(true),
				Log:      logr.Discard(),
			}
			r := &testReconciler{updateErr: tt.updateErr, statusUpdateErr: tt.statusUpdateErr}

			if _, err := registry.HandleExecution(r, req); (err != nil) != tt.wantErr {
				t.Fatalf("HandleExecution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
