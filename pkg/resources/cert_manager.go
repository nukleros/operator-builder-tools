package resources

import (
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	IssuerKind        = cmv1.IssuerKind
	ClusterIssuerKind = cmv1.ClusterIssuerKind
	CertificateKind   = cmv1.CertificateKind
)

// IssuerResource represents a cert-manager Issuer object.
type IssuerResource struct {
	Object cmv1.Issuer
}

// ClusterIssuerResource represents a cert-manager ClusterIssuer object.
type ClusterIssuerResource struct {
	Object cmv1.ClusterIssuer
}

// CertificateResource represents a cert-manager Certificate object.
type CertificateResource struct {
	Object cmv1.Certificate
}

// NewIssuerResource creates and returns a new IssuerResource.
func NewIssuerResource(object client.Object) (*IssuerResource, error) {
	issuer := &cmv1.Issuer{}

	err := ToTyped(issuer, object)
	if err != nil {
		return nil, err
	}

	return &IssuerResource{Object: *issuer}, nil
}

// NewClusterIssuerResource creates and returns a new ClusterIssuerResource.
func NewClusterIssuerResource(object client.Object) (*ClusterIssuerResource, error) {
	clusterIssuer := &cmv1.ClusterIssuer{}

	err := ToTyped(clusterIssuer, object)
	if err != nil {
		return nil, err
	}

	return &ClusterIssuerResource{Object: *clusterIssuer}, nil
}

// NewCertificateResource creates and returns a new CertificateResource.
func NewCertificateResource(object client.Object) (*CertificateResource, error) {
	cert := &cmv1.Certificate{}

	err := ToTyped(cert, object)
	if err != nil {
		return nil, err
	}

	return &CertificateResource{Object: *cert}, nil
}

// IsReady checks to see if an Issuer is ready.
func (issuer *IssuerResource) IsReady() (bool, error) {
	return issuerIsReady(issuer.Object.Status.Conditions)
}

// IsReady checks to see if an Issuer is ready.
func (clusterIssuer *ClusterIssuerResource) IsReady() (bool, error) {
	return issuerIsReady(clusterIssuer.Object.Status.Conditions)
}

// IsReady checks to see if an Certificate is ready.
func (cert *CertificateResource) IsReady() (bool, error) {
	for _, condition := range cert.Object.Status.Conditions {
		if condition.Type == cmv1.CertificateConditionReady && condition.Status == cmmetav1.ConditionTrue {
			return true, nil
		}
	}

	return false, nil
}

// issuerIsReady determines if either an Issuer or a ClusterIssuer resource is ready.
func issuerIsReady(conditions []cmv1.IssuerCondition) (bool, error) {
	for _, condition := range conditions {
		if condition.Type == cmv1.IssuerConditionReady && condition.Status == cmmetav1.ConditionTrue {
			return true, nil
		}
	}

	return false, nil
}
