/*
	SPDX-License-Identifier: MIT
*/

package resources

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNamespaceResource_IsReady(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			namespace := &NamespaceResource{
				Namespace: tt.fields.Namespace,
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
