/*
	SPDX-License-Identifier: MIT
*/

package resources_test

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/nukleros/operator-builder-tools/pkg/resources"
)

func TestNamespaceResource_IsReady(t *testing.T) {
	t.Parallel()

	type fields struct {
		Namespace v1.Namespace
	}

	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name:    "namespace should be ready",
			want:    true,
			wantErr: false,
			fields: fields{
				Namespace: v1.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-namespace",
						Namespace: "",
					},
					Status: v1.NamespaceStatus{
						Phase: v1.NamespaceActive,
					},
				},
			},
		},
		{
			name:    "namespace should be ready (with namespace)",
			want:    true,
			wantErr: false,
			fields: fields{
				Namespace: v1.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ready-namespace",
						Namespace: "unknown-namespace",
					},
					Status: v1.NamespaceStatus{
						Phase: v1.NamespaceActive,
					},
				},
			},
		},
		{
			name:    "namespace should not be ready (terminating)",
			want:    false,
			wantErr: false,
			fields: fields{
				Namespace: v1.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "terminating-namespace",
						Namespace: "",
					},
					Status: v1.NamespaceStatus{
						Phase: v1.NamespaceTerminating,
					},
				},
			},
		},
		{
			name:    "namespace should not be ready (unknown status)",
			want:    false,
			wantErr: false,
			fields: fields{
				Namespace: v1.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "terminating-namespace",
						Namespace: "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			namespace := &resources.NamespaceResource{
				Object: tt.fields.Namespace,
			}
			got, err := namespace.IsReady()
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceResource.IsReady() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("NamespaceResource.IsReady() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNamespaceResource(t *testing.T) {
	t.Parallel()

	type args struct {
		object client.Object
	}

	tests := []struct {
		name    string
		args    args
		want    *resources.NamespaceResource
		wantErr bool
	}{
		{
			name: "namespace should be created",
			want: &resources.NamespaceResource{
				Object: v1.Namespace{},
			},
			wantErr: false,
			args: args{
				object: &v1.Namespace{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resources.NewNamespaceResource(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNamespaceResource() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNamespaceResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
