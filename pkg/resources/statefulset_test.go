/*
	SPDX-License-Identifier: MIT
*/

package resources_test

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

func TestNewStatefulSetResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.StatefulSetResource
		wantErr bool
	}{
		{
			name: "statefulset should be created",
			want: &resources.StatefulSetResource{
				Object: appsv1.StatefulSet{},
			},
			wantErr: false,
			args: args{
				object: &appsv1.StatefulSet{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewStatefulSetResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStatefulSetResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatefulSetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
