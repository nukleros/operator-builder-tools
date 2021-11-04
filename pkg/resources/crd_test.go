/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"reflect"
	"testing"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCRDResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *CRDResource
		wantErr bool
	}{
		{
			name: "crd should be created",
			want: &CRDResource{
				Object: extensionsv1.CustomResourceDefinition{},
			},
			wantErr: false,
			args: args{
				object: &extensionsv1.CustomResourceDefinition{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCRDResource(tt.args.object)
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
