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

func TestNewDaemonSetResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *DaemonSetResource
		wantErr bool
	}{
		{
			name: "daemonset should be created",
			want: &DaemonSetResource{
				Object: appsv1.DaemonSet{},
			},
			wantErr: false,
			args: args{
				object: &appsv1.DaemonSet{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDaemonSetResource(tt.args.object)
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
