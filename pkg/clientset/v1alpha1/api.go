package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

type SecurityV1Alpha1Interface interface {
	Audits(namespace string) AuditInterface
	Whitelists(namespace string) WhiteListInterface
}

type SecurityV1Alpha1Client struct {
	restClient rest.Interface
}

func SecurityConfig(c *rest.Config) (*SecurityV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &SecurityV1Alpha1Client{restClient: client}, nil
}

func (c *SecurityV1Alpha1Client) Audits(namespace string) AuditInterface {
	return &auditClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *SecurityV1Alpha1Client) Whitelists(namespace string) WhiteListInterface {
	return &whitelistClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
