/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"reflect"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewJobResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *JobResource
		wantErr bool
	}{
		{
			name: "job should be created",
			want: &JobResource{
				Object: batchv1.Job{},
			},
			wantErr: false,
			args: args{
				object: &batchv1.Job{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJobResource(tt.args.object)
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
