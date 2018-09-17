package whitelist

import (
	"strings"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

func getWhiteList() ([]v1alpha1.WhiteList, error) {
	var config *rest.Config
	var err error

	glog.Info("using in-cluster configuration")
	config, err = rest.InClusterConfig()

	v1alpha1.AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		glog.Error(err)
	}

	result := v1alpha1.WhiteListList{}
	err = exampleRestClient.Get().Resource("whitelists").Do().Into(&result)

	return result.Items, err
}

func CheckWhiteList(l map[string]string, s string) bool {
	wl, err := getWhiteList()
	if err != nil {
		glog.Errorf("Reading whitelists failed: %s ", err)
	}
	release := l["release"]
	if release != "" {
		for _, res := range wl {
			if release == res.Spec.ReleaseName {
				return true
			}
		}
	} else {
		for _, res := range wl {
			if strings.Contains(s, res.Spec.ReleaseName) {
				return true
			}
		}
	}
	return false
}
