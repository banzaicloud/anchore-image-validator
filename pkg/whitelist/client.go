package whitelist

import (
	"log"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

func GetWhiteList() ([]v1alpha1.WhiteList, error) {
	var config *rest.Config
	var err error

	log.Printf("using in-cluster configuration")
	config, err = rest.InClusterConfig()

	v1alpha1.AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	result := v1alpha1.WhiteListList{}
	err = exampleRestClient.Get().Resource("whitelists").Do().Into(&result)

	return result.Items, err
}

func CheckWhiteList(s string) bool {
	wl, err := GetWhiteList()
	if err != nil {
		panic(err)
	}

	for _, res := range wl {
		if strings.Contains(s, res.Spec.ReleaseName) {
			return true
		}
	}
	return false
}
