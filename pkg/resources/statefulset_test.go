/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewStatefulSetResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *StatefulSetResource
		wantErr bool
	}{
		{
			name: "statefulset should be created",
			want: &StatefulSetResource{
				Object: appsv1.StatefulSet{},
			},
			wantErr: false,
			args: args{
				object: &appsv1.StatefulSet{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStatefulSetResource(tt.args.object)
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
