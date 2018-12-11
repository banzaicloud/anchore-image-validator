// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/banzaicloud/anchore-image-validator/pkg/anchore"
	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	clientV1alpha1 "github.com/banzaicloud/anchore-image-validator/pkg/clientset/v1alpha1"
	"github.com/openshift/generic-admission-server/pkg/cmd"
	"github.com/sirupsen/logrus"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var securityClientSet *clientV1alpha1.SecurityV1Alpha1Client

type admissionHook struct {
	reservationClient dynamic.ResourceInterface
	lock              sync.RWMutex
	initialized       bool
}

const apiServiceResource = "imagechecks"

var (
	apiServiceGroup     = os.Getenv("ANCHORE_APISERVICE_GROUP")
	apiServiceVersion   = os.Getenv("ANCHORE_APISERVICE_VERSION")
	anchoreReleaseName  = os.Getenv("ANCHORE_RELEASE_NAME")
	kubernetesNameSpace = os.Getenv("KUBERNETES_NAMESPACE")
)

func main() {
	var config *rest.Config
	var err error

	config, err = rest.InClusterConfig()
	if err != nil {
		logrus.Error(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)
	securityClientSet, err = clientV1alpha1.SecurityConfig(config)
	if err != nil {
		logrus.Error(err)
	}

	installValidatingWebhookConfig(config)

	cmd.RunAdmissionServer(&admissionHook{})
}

func (a *admissionHook) ValidatingResource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    apiServiceGroup,
			Version:  apiServiceVersion,
			Resource: apiServiceResource,
		},
		"imagecheck"
}

func (a *admissionHook) Validate(admissionSpec *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	status := &admissionv1beta1.AdmissionResponse{
		Allowed: true,
		UID:     admissionSpec.UID,
		Result:  &metav1.Status{Status: "Success", Message: ""}}

	if admissionSpec.Kind.Kind == "Pod" {
		whitelists, err := securityClientSet.Whitelists("default").List(metav1.ListOptions{})
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.WithFields(logrus.Fields{
				"whitelists": whitelists.Items,
			}).Debug("Whitelists found")
		}
		pod := v1.Pod{}
		json.Unmarshal(admissionSpec.Object.Raw, &pod)
		logrus.WithFields(logrus.Fields{
			"PodName":    pod.Name,
			"NameSpace":  pod.Namespace,
			"Labels":     pod.Labels,
			"Anotations": pod.Annotations,
		}).Debug("Pod details")

		var result []string
		var message string
		var auditImages []v1alpha1.AuditImage
		r, f := getReleaseName(pod.Labels, pod.Name)
		for _, container := range pod.Spec.Containers {
			image := container.Image
			logrus.WithFields(logrus.Fields{
				"image": image,
			}).Info("Checking image")
			auditImage, ok := anchore.CheckImage(image)
			if !ok {
				status.Result.Status = "Failure"
				status.Allowed = false
				if checkWhiteList(whitelists.Items, r, f) {
					status.Result.Status = "Success"
					status.Allowed = true
					logrus.WithFields(logrus.Fields{
						"PodName": pod.Name,
					}).Info("Whitelisted release")
				}
				message = fmt.Sprintf("Image failed policy check: %s", image)
				status.Result.Message = message
				logrus.WithFields(logrus.Fields{
					"image": image,
				}).Warning("Image failed policy check")
			} else {
				message = fmt.Sprintf("Image passed policy check: %s", image)
				logrus.WithFields(logrus.Fields{
					"image": image,
				}).Warning("Image passed policy check")
			}
			result = append(result, message)
			auditImages = append(auditImages, auditImage)
		}

		fr := "false"
		if f {
			fr = "true"
		}
		action := "reject"
		if status.Allowed {
			action = "allowed"
		}
		owners := pod.GetOwnerReferences()
		var auditName string
		if len(owners) > 0 {
			auditName = strings.ToLower(owners[0].Kind) + "-" + strings.ToLower(owners[0].Name)
		} else {
			auditName = pod.Name
		}

		ainfo := auditInfo{
			name:        auditName,
			labels:      map[string]string{"fakerelease": fr},
			releaseName: r,
			resource:    "Pod",
			images:      auditImages,
			result:      result,
			action:      action,
			state:       "",
			owners:      owners,
		}

		createOrUpdateAudit(ainfo)
		logrus.WithFields(logrus.Fields{
			"Status": status,
		}).Debug("Security scan status")
	}
	return status
}

func (a *admissionHook) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
