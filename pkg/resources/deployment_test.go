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

func TestDeploymentResource_IsReady(t *testing.T) {
	var randomInt int32
	randomInt = 1

	type fields struct {
		parent *appsv1.Deployment
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name:    "deployment should be ready",
			want:    true,
			wantErr: false,
			fields: fields{
				parent: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "ready",
						Namespace:  "ready",
						Generation: int64(randomInt),
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &randomInt,
					},
					Status: appsv1.DeploymentStatus{
						Replicas:           randomInt,
						ReadyReplicas:      randomInt,
						ObservedGeneration: int64(randomInt),
					},
				},
			},
		},
		{
			name:    "deployment should not be ready (replicas)",
			want:    false,
			wantErr: false,
			fields: fields{
				parent: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "not-ready-replicas",
						Namespace:  "not-ready-replicas",
						Generation: int64(randomInt),
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &randomInt,
					},
					Status: appsv1.DeploymentStatus{
						Replicas:           randomInt,
						ReadyReplicas:      randomInt + 1,
						ObservedGeneration: int64(randomInt),
					},
				},
			},
		},
		{
			name:    "deployment should not be ready (name)",
			want:    false,
			wantErr: false,
			fields: fields{
				parent: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:  "not-ready-name",
						Generation: int64(randomInt),
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &randomInt,
					},
					Status: appsv1.DeploymentStatus{
						Replicas:           randomInt,
						ReadyReplicas:      randomInt,
						ObservedGeneration: int64(randomInt),
					},
				},
			},
		},
		{
			name:    "deployment should not be ready (unavailable replicas)",
			want:    false,
			wantErr: false,
			fields: fields{
				parent: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-ready-unavailable",
						Namespace: "not-ready-unavailable",
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &randomInt,
					},
					Status: appsv1.DeploymentStatus{
						Replicas:            randomInt,
						ReadyReplicas:       randomInt,
						UnavailableReplicas: 1,
					},
				},
			},
		},
		{
			name:    "deployment should not be ready (empty)",
			want:    false,
			wantErr: false,
			fields: fields{
				parent: &appsv1.Deployment{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployment := &DeploymentResource{
				*tt.fields.parent,
			}
			got, err := deployment.IsReady()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentResource.IsReady() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeploymentResource.IsReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeploymentResource(t *testing.T) {
	type args struct {
		object metav1.Object
	}
	tests := []struct {
		name    string
		args    args
		want    *DeploymentResource
		wantErr bool
	}{
		{
			name: "deployment should be created",
			want: &DeploymentResource{
				Object: appsv1.Deployment{},
			},
			wantErr: false,
			args: args{
				object: &appsv1.Deployment{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDeploymentResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDeploymentResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeploymentResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
