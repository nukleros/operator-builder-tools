/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewSecretResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *SecretResource
		wantErr bool
	}{
		{
			name: "secret should be created",
			want: &SecretResource{
				Object: v1.Secret{},
			},
			wantErr: false,
			args: args{
				object: &v1.Secret{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSecretResource(tt.args.object)
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
