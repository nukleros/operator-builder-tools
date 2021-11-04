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

func TestNewConfigMapResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *ConfigMapResource
		wantErr bool
	}{
		{
			name: "configmap should be created",
			want: &ConfigMapResource{
				Object: v1.ConfigMap{},
			},
			wantErr: false,
			args: args{
				object: &v1.ConfigMap{},
			},
		},
		{
			name:    "configmap should not be created",
			want:    nil,
			wantErr: true,
			args: args{
				object: &v1.Namespace{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigMapResource(tt.args.object)
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
