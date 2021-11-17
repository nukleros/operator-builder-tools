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

func TestNewDaemonSetResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.DaemonSetResource
		wantErr bool
	}{
		{
			name: "daemonset should be created",
			want: &resources.DaemonSetResource{
				Object: appsv1.DaemonSet{},
			},
			wantErr: false,
			args: args{
				object: &appsv1.DaemonSet{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewDaemonSetResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDaemonSetResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDaemonSetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
