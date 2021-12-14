/*
	SPDX-License-Identifier: MIT
*/

package resources_test

import (
	"reflect"
	"testing"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

func TestNewCRDResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.CRDResource
		wantErr bool
	}{
		{
			name: "crd should be created",
			want: &resources.CRDResource{
				Object: extensionsv1.CustomResourceDefinition{},
			},
			wantErr: false,
			args: args{
				object: &extensionsv1.CustomResourceDefinition{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewCRDResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCRDResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCRDResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
