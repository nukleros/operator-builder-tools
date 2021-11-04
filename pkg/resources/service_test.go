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

func TestNewServiceResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *ServiceResource
		wantErr bool
	}{
		{
			name: "service should be created",
			want: &ServiceResource{
				Object: v1.Service{},
			},
			wantErr: false,
			args: args{
				object: &v1.Service{},
			},
		},
		{
			name:    "service should not be created",
			want:    nil,
			wantErr: true,
			args: args{
				object: &v1.Namespace{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServiceResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServiceResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServiceResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
