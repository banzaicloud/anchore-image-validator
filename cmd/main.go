package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"github.com/openshift/generic-admission-server/pkg/cmd"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/banzaicloud/anchore-image-validator/pkg/anchore"
	"github.com/banzaicloud/anchore-image-validator/pkg/whitelist"
)

var log logrus.New()

type admissionHook struct {
	reservationClient dynamic.ResourceInterface
	lock              sync.RWMutex
	initialized       bool
}

func main() {
	cmd.RunAdmissionServer(&admissionHook{})
}

func (a *admissionHook) ValidatingResource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.anchore.io",
			Version:  "v1beta1",
			Resource: "imagechecks",
		},
		"imagecheck"
}

func (a *admissionHook) Validate(admissionSpec *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	status := &admissionv1beta1.AdmissionResponse{
		Allowed: true,
		UID:     admissionSpec.UID,
		Result:  &metav1.Status{Status: "Success", Message: ""}}

	if admissionSpec.Kind.Kind == "Pod" {
		pod := v1.Pod{}
		json.Unmarshal(admissionSpec.Object.Raw, &pod)
		glog.Infof("\n --- DEBUG MESSAGE NAME: --- %s \n", pod.Name)
		glog.Infof("\n --- DEBUG MESSAGE NAMESPACE: --- %s \n", pod.Namespace)
		glog.Infof("\n --- DEBUG MESSAGE LABELS: --- %s \n", pod.Labels)
		glog.Infof("\n --- DEBUG MESSAGE ANNOTATIONS: --- %s \n", pod.Annotations)
		for _, container := range pod.Spec.Containers {
			image := container.Image
			glog.Infof("Checking image: %s", image)
			if !anchore.CheckImage(image) {
				status.Result.Status = "Failure"
				status.Allowed = false
				if whitelist.CheckWhiteList(pod.Labels, pod.Name) {
					status.Result.Status = "Success"
					status.Allowed = true
					glog.Infof("Whitelisted release in case of pod: %s: ", pod.Name)
				}
				message := fmt.Sprintf("Image failed policy check: %s", image)
				status.Result.Message = message
				glog.Warning(message)
				glog.Infof("\n --- STATUS ---: %s\n", status)
				return status
			} else {
				glog.Infof("Image passed policy check: %s", image)
			}
		}
	}
	glog.Infof("\n --- STATUS ---: %s\n", status)
	glog.Flush()
	return status
}

func (a *admissionHook) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
