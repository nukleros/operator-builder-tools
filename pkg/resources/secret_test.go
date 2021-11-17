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

func TestNewSecretResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.SecretResource
		wantErr bool
	}{
		{
			name: "secret should be created",
			want: &resources.SecretResource{
				Object: v1.Secret{},
			},
			wantErr: false,
			args: args{
				object: &v1.Secret{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewSecretResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSecretResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecretResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
