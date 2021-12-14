/*
	SPDX-License-Identifier: MIT
*/

package resources_test

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

func TestNewConfigMapResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.ConfigMapResource
		wantErr bool
	}{
		{
			name: "configmap should be created",
			want: &resources.ConfigMapResource{
				Object: v1.ConfigMap{},
			},
			wantErr: false,
			args: args{
				object: &v1.ConfigMap{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewConfigMapResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigMapResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigMapResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
