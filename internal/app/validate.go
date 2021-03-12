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
	"time"

	"github.com/banzaicloud/anchore-image-validator/pkg/anchore"
	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	"github.com/dgraph-io/ristretto"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"logur.dev/logur"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func validate(ar *admissionv1beta1.AdmissionReview,
	logger logur.Logger,
	c client.Client,
	cache *ristretto.Cache,
	cacheTTL time.Duration) *admissionv1beta1.AdmissionResponse {
	req := ar.Request

	logger.Info("AdmissionReview for", map[string]interface{}{
		"Kind":      req.Kind,
		"Namespsce": req.Namespace,
		"Resource":  req.Resource,
		"UserInfo":  req.UserInfo,
	})

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

		r, f := getReleaseName(pod.Labels, pod.Name)

		if checkWhiteList(whitelists.Items, r, f) {
			logger.Info("Whitelisted release", map[string]interface{}{
				"PodName": pod.Name,
			})

			go checkImage(&pod, logger, c, cache, cacheTTL, true)

			return &admissionv1beta1.AdmissionResponse{
				Allowed: true,
				Result: &metav1.Status{
					Status:  "Success",
					Reason:  "",
					Message: "Whitelisted release",
				},
			}
		}

		return checkImage(&pod, logger, c, cache, cacheTTL, false)
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
	logger logur.Logger,
	c client.Client,
	cache *ristretto.Cache,
	cacheTTL time.Duration,
	isWhiteListed bool) *admissionv1beta1.AdmissionResponse {
	result := []string{}
	auditImages := []v1alpha1.AuditImage{}

	var message string

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

		isCached := false

		logger.Debug("Checking cache", map[string]interface{}{
			"PodName": pod.Name,
		})

		_, found := cache.Get(image)
		if found {
			logger.Debug("not submitting for scan again within ttl", map[string]interface{}{
				"image": image,
			})

			isCached = true
		} else {
			logger.Debug("submitting for scan the first time within ttl", map[string]interface{}{
				"image": image,
			})
			cache.SetWithTTL(image, "submitted", 100, cacheTTL)
		}

		auditImage, ok := anchore.CheckImage(image, isCached, isWhiteListed)

		if !ok {
			resp.Result.Status = "Failure"
			resp.Allowed = false
			message = fmt.Sprintf("Image failed policy check: %s", image)
			resp.Result.Message = message

			logger.Warn("Image failed policy check", map[string]interface{}{
				"image": image,
			})
		} else {
			if isWhiteListed {
				message = fmt.Sprintf("Whitelisted release: %s", image)
				logger.Info("skipping policy enforcement for whitelisted release", map[string]interface{}{
					"image": image,
				})
			} else {
				message = fmt.Sprintf("Image passed policy check: %s", image)
				logger.Warn("Image passed policy check", map[string]interface{}{
					"image": image,
				})
			}
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
