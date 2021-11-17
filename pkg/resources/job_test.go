/*
	SPDX-License-Identifier: MIT
*/

package resources_test

import (
	"reflect"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

func TestNewJobResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.JobResource
		wantErr bool
	}{
		{
			name: "job should be created",
			want: &resources.JobResource{
				Object: batchv1.Job{},
			},
			wantErr: false,
			args: args{
				object: &batchv1.Job{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewJobResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJobResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJobResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
