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

func TestServiceResource_IsReady(t *testing.T) {
	type fields struct {
		Object v1.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name:    "service should be ready (external)",
			want:    true,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type: v1.ServiceTypeExternalName,
					},
				},
			},
		},
		{
			name:    "service should be ready (clusterip string)",
			want:    true,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type:      v1.ServiceTypeClusterIP,
						ClusterIP: "1.1.1.1",
					},
				},
			},
		},
		{
			name:    "service should be ready (clusterip slice)",
			want:    true,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type:       v1.ServiceTypeClusterIP,
						ClusterIPs: []string{"1.1.1.1"},
					},
				},
			},
		},
		{
			name:    "service should be ready (load balancer)",
			want:    true,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type: v1.ServiceTypeLoadBalancer,
					},
					Status: v1.ServiceStatus{
						LoadBalancer: v1.LoadBalancerStatus{
							Ingress: []v1.LoadBalancerIngress{
								{
									IP: "1.1.1.1",
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "service should not be ready (clusterip string)",
			want:    false,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type:      v1.ServiceTypeClusterIP,
						ClusterIP: "",
					},
				},
			},
		},
		{
			name:    "service should be ready (clusterip slice)",
			want:    false,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type:       v1.ServiceTypeClusterIP,
						ClusterIPs: []string{},
					},
				},
			},
		},
		{
			name:    "service should not be ready (load balancer)",
			want:    false,
			wantErr: false,
			fields: fields{
				Object: v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-service",
						Namespace: "ready-namespace",
					},
					Spec: v1.ServiceSpec{
						Type: v1.ServiceTypeLoadBalancer,
					},
					Status: v1.ServiceStatus{
						LoadBalancer: v1.LoadBalancerStatus{
							Ingress: []v1.LoadBalancerIngress{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ServiceResource{
				Object: tt.fields.Object,
			}
			got, err := service.IsReady()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceResource.IsReady() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceResource.IsReady() = %v, want %v", got, tt.want)
			}
		})
	}
}
