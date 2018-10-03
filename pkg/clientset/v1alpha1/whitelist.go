package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

// WhiteListInterface for whitelist
type WhiteListInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.WhiteList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.WhiteListItem, error)
	Create(*v1alpha1.WhiteListItem) (*v1alpha1.WhiteListItem, error)
	Delete(name string, options *metav1.DeleteOptions) error
}

type whitelistClient struct {
	restClient rest.Interface
	ns         string
}

func (c *whitelistClient) List(opts metav1.ListOptions) (*v1alpha1.WhiteList, error) {
	result := v1alpha1.WhiteList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("whitelists").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *whitelistClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.WhiteListItem, error) {
	result := v1alpha1.WhiteListItem{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("whitelists").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *whitelistClient) Create(whiteListItem *v1alpha1.WhiteListItem) (*v1alpha1.WhiteListItem, error) {
	result := v1alpha1.WhiteListItem{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("whitelists").
		Body(whiteListItem).
		Do().
		Into(&result)

	return &result, err
}

func (c *whitelistClient) Delete(name string, options *metav1.DeleteOptions) error {

	return c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("whitelists").
		Name(name).
		Body(options).
		Do().
		Error()
}
