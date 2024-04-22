package resources

import (
	"fmt"

	"github.com/nukleros/operator-builder-tools/pkg/controller/workload"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MutatingWebhookConfigurationKind    = "MutatingWebhookConfiguration"
	MutatingWebhookConfigurationVersion = "admissionregistration.k8s.io/v1"

	ValidatingWebhookConfigurationKind    = "ValidatingWebhookConfiguration"
	ValidatingWebhookConfigurationVersion = "admissionregistration.k8s.io/v1"
)

// MutatingWebhookConfigurationResource represents a Kubernetes
// MutatingWebhookConfiguration object.
type MutatingWebhookConfigurationResource struct {
	Object     admissionv1.MutatingWebhookConfiguration
	Reconciler workload.Reconciler
	Request    *workload.Request
}

// ValidatingWebhookConfigurationResource represents a Kubernetes
// ValidatingWebhookConfiguration object.
type ValidatingWebhookConfigurationResource struct {
	Object     admissionv1.ValidatingWebhookConfiguration
	Reconciler workload.Reconciler
	Request    *workload.Request
}

// NewMutatingWebhookConfigurationResource creates and returns a new MutatingWebhookConfigurationResource.
func NewMutatingWebhookConfigurationResource(
	r workload.Reconciler,
	req *workload.Request,
	object client.Object,
) (*MutatingWebhookConfigurationResource, error) {
	webhook := &admissionv1.MutatingWebhookConfiguration{}

	err := ToTyped(webhook, object)
	if err != nil {
		return nil, err
	}

	return &MutatingWebhookConfigurationResource{
		Object:     *webhook,
		Reconciler: r,
		Request:    req,
	}, nil
}

// NewMValidatingWebhookConfigurationResource creates and returns a new ValidatingWebhookConfigurationResource.
func NewValidatingWebhookConfigurationResource(
	r workload.Reconciler,
	req *workload.Request,
	object client.Object,
) (*ValidatingWebhookConfigurationResource, error) {
	webhook := &admissionv1.ValidatingWebhookConfiguration{}

	err := ToTyped(webhook, object)
	if err != nil {
		return nil, err
	}

	return &ValidatingWebhookConfigurationResource{
		Object:     *webhook,
		Reconciler: r,
		Request:    req,
	}, nil
}

// IsReady performs the logic to determine if a MutatingWebhookConfiguration is ready.
func (webhook *MutatingWebhookConfigurationResource) IsReady() (bool, error) {
	return isReady(webhook, webhook.Reconciler, webhook.Request)
}

// IsReady performs the logic to determine if a MutatingWebhookConfiguration is ready.
func (webhook *ValidatingWebhookConfigurationResource) IsReady() (bool, error) {
	return isReady(webhook, webhook.Reconciler, webhook.Request)
}

// GetServiceStubs gets the service stubs objects from a MutatingWebhookConfigurationResource.  The stubs
// are used to lookup the underlying services associated with the webhook.
func (webhook *MutatingWebhookConfigurationResource) GetServiceStubs() []v1.Service {
	services := []v1.Service{}

	for _, mutating := range webhook.Object.Webhooks {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mutating.ClientConfig.Service.Name,
				Namespace: mutating.ClientConfig.Service.Namespace,
			},
		}

		services = append(services, service)
	}

	return services
}

// GetServiceStubs gets the service stubs objects from a ValidatingWebhookConfigurationResource.  The stubs
// are used to lookup the underlying services associated with the webhook.
func (webhook *ValidatingWebhookConfigurationResource) GetServiceStubs() []v1.Service {
	services := []v1.Service{}

	for _, validating := range webhook.Object.Webhooks {
		service := v1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       ServiceKind,
				APIVersion: ServiceVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      validating.ClientConfig.Service.Name,
				Namespace: validating.ClientConfig.Service.Namespace,
			},
		}

		services = append(services, service)
	}

	return services
}

// isReady determines if a webhook resource is ready.
func isReady(serviceRetriever ServiceStubRetriever, r workload.Reconciler, req *workload.Request) (bool, error) {
	for _, service := range serviceRetriever.GetServiceStubs() {
		service, err := NewServiceResource(&service)
		if err != nil {
			return false, err
		}

		endpoints, err := service.GetEndpoints(r, req)
		if err != nil {
			return false, fmt.Errorf("unable to retrieve endpoints from service - %w", err)
		}

		ready, err := endpoints.IsReady()
		if err != nil {
			return false, fmt.Errorf("unable to determine endpoint readiness = %w", err)
		}

		if !ready {
			return false, nil
		}
	}

	return true, nil
}
