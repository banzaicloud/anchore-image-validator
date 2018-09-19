package whitelist

import (
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
)

func getWhiteList() ([]v1alpha1.WhiteList, error) {
	var config *rest.Config
	var err error

	logrus.Info("using in-cluster configuration")
	config, err = rest.InClusterConfig()

	v1alpha1.AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		logrus.Error(err)
	}

	result := v1alpha1.WhiteListList{}
	err = exampleRestClient.Get().Resource("whitelists").Do().Into(&result)

	return result.Items, err
}

func CheckWhiteList(l map[string]string, s string) bool {
	wl, err := getWhiteList()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Reading whitelists failed")
	}
	release := l["release"]
	if release != "" {
		logrus.WithFields(logrus.Fields{
			"release": release,
		}).Info("Check whitelist")
		for _, res := range wl {
			if release == res.Spec.ReleaseName {
				return true
			}
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"PodName": s,
		}).Info("Missing release label, using PodName")
		for _, res := range wl {
			fakeRelease := string(res.Spec.ReleaseName + "-")
			if strings.Contains(s, fakeRelease) {
				return true
			}
		}
	}
	return false
}
