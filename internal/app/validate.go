// Copyright © 2019 Banzai Cloud.
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

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/banzaicloud/anchore-image-validator/pkg/anchore"
	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"logur.dev/logur"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func validate(ar *admissionv1beta1.AdmissionReview,
	logger logur.Logger, c client.Client) *admissionv1beta1.AdmissionResponse {
	req := ar.Request

	logger.Info("AdmissionReview for", map[string]interface{}{
		"Kind":      req.Kind,
		"Namespsce": req.Namespace,
		"Resource":  req.Resource,
		"UserInfo":  req.UserInfo})

	if req.Kind.Kind == "Pod" {
		whitelists := &v1alpha1.WhiteListItemList{}

		if err := c.List(context.Background(), whitelists); err != nil {
			logger.Error("cannot list whitelistimets", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			logger.Debug("whitelists found", map[string]interface{}{
				"whitelists": whitelists.Items,
			})
		}

		pod := v1.Pod{}
		if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
			logger.Error("could not unmarshal raw object")

			return &admissionv1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}

		return checkImage(&pod, whitelists, logger, c)
	}

	return &admissionv1beta1.AdmissionResponse{
		Allowed: true,
		Result: &metav1.Status{
			Status:  "Success",
			Reason:  "",
			Message: "",
		},
	}
}

func checkImage(pod *v1.Pod,
	wl *v1alpha1.WhiteListItemList,
	logger logur.Logger,
	c client.Client) *admissionv1beta1.AdmissionResponse {
	result := []string{}
	auditImages := []v1alpha1.AuditImage{}
	message := ""

	resp := &admissionv1beta1.AdmissionResponse{
		Allowed: true,
		Result: &metav1.Status{
			Status:  "Success",
			Reason:  "",
			Message: "",
		},
	}

	r, f := getReleaseName(pod.Labels, pod.Name)

	for _, container := range pod.Spec.Containers {
		image := container.Image

		logger.Debug("Checking image", map[string]interface{}{
			"image": image,
		})

		auditImage, ok := anchore.CheckImage(image)

		if !ok {
			resp.Result.Status = "Failure"
			resp.Allowed = false

			if checkWhiteList(wl.Items, r, f) {
				resp.Result.Status = "Success"
				resp.Allowed = true

				logger.Info("Whitelisted release", map[string]interface{}{
					"PodName": pod.Name,
				})
			}
			message = fmt.Sprintf("Image failed policy check: %s", image)
			resp.Result.Message = message

			logger.Warn("Image failed policy check", map[string]interface{}{
				"image": image,
			})
		} else {
			message = fmt.Sprintf("Image passed policy check: %s", image)

			logger.Warn("Image passed policy check", map[string]interface{}{
				"image": image,
			})
		}

		result = append(result, message)
		auditImages = append(auditImages, auditImage)
	}

	fr := "false"
	if f {
		fr = "true"
	}

	action := "reject"
	if resp.Allowed {
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

	createOrUpdateAudit(ainfo, c)
	logger.Debug("Security scan status", map[string]interface{}{
		"Status": resp,
	})

	return resp
}
