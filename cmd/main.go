package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/golang/glog"
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
		if !whitelist.CheckWhiteList(pod.Name) {
			for _, container := range pod.Spec.Containers {
				image := container.Image
				glog.Infof("Checking image: %s", image)
				if !anchore.CheckImage(image) {
					status.Result.Status = "Failure"
					status.Allowed = false
					message := fmt.Sprintf("Image failed policy check: %s", image)
					status.Result.Message = message
					glog.Warning(message)
					return status
				} else {
					glog.Infof("Image passed policy check: %s", image)
				}
			}
		} else {
			glog.Info("Whitelisted image name, policy check skipped: " + pod.Name)
		}
	}
	glog.Flush()
	return status
}

func (a *admissionHook) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
