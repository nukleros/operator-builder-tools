/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
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
	type args struct {
		resource *Resource
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "deployment should be ready",
			want:    true,
			wantErr: false,
			args: args{
				resource: &Resource{
					ResourceCommon: ResourceCommon{
						Name:      "ready",
						Namespace: "ready",
					},
				},
			},
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
			args: args{
				resource: &Resource{
					ResourceCommon: ResourceCommon{
						Name:      "not-ready-replicas",
						Namespace: "not-ready-replicas",
					},
				},
			},
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
			args: args{
				resource: &Resource{
					ResourceCommon: ResourceCommon{
						Name:      "not-ready-name",
						Namespace: "not-ready-name",
					},
				},
			},
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
			name:    "deployment should not be ready (observed generation)",
			want:    false,
			wantErr: false,
			args: args{
				resource: &Resource{
					ResourceCommon: ResourceCommon{
						Name:      "not-ready-generation",
						Namespace: "not-ready-generation",
					},
				},
			},
			fields: fields{
				parent: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "not-ready-generation",
						Namespace:  "not-ready-generation",
						Generation: int64(randomInt),
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &randomInt,
					},
					Status: appsv1.DeploymentStatus{
						Replicas:           randomInt,
						ReadyReplicas:      randomInt,
						ObservedGeneration: int64(randomInt + 1),
					},
				},
			},
		},
		{
			name:    "deployment should not be ready (empty)",
			want:    false,
			wantErr: false,
			args:    args{},
			fields: fields{
				parent: &appsv1.Deployment{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployment := &DeploymentResource{
				parent: tt.fields.parent,
			}
			got, err := deployment.IsReady(tt.args.resource)
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
