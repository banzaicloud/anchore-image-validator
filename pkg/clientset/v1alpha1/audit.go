package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

// AuditInterface for audit
type AuditInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.AuditList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Audit, error)
	Create(*v1alpha1.Audit) (*v1alpha1.Audit, error)
}

type auditClient struct {
	restClient rest.Interface
	ns         string
}

func (c *auditClient) List(opts metav1.ListOptions) (*v1alpha1.AuditList, error) {
	result := v1alpha1.AuditList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("audits").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *auditClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.Audit, error) {
	result := v1alpha1.Audit{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("audits").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *auditClient) Create(audit *v1alpha1.Audit) (*v1alpha1.Audit, error) {
	result := v1alpha1.Audit{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("audits").
		Body(audit).
		Do().
		Into(&result)

	return &result, err
}
