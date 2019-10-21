/*
Copyright 2019 Banzai Cloud.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"encoding/json"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"logur.dev/logur"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func validate(ar *admissionv1beta1.AdmissionReview, logger logur.Logger, c client.Client) *admissionv1beta1.AdmissionResponse {
	req := ar.Request
	logger.Info("AdmissionReview for", map[string]interface{}{
		"Kind":      req.Kind,
		"Namespsce": req.Namespace,
		"Resource":  req.Resource,
		"UserInfo":  req.UserInfo})

	switch req.Kind.Kind {
	case "Pod":
		pod := v1.Pod{}
		if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
			logger.Error("could not unmarshal raw object")
			return &admissionv1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}

		ok, err := checkImage(&pod, pod.GetNamespace(), logger)
		if err != nil {
			return &admissionv1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Reason: metav1.StatusReason(err.Error()),
				},
			}
		}
		if !ok {
			return &admissionv1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Reason: "scan results are above treshold",
				},
			}
		}
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

func checkImage(pod *v1.Pod, namespave string, logger logur.Logger) (bool, error) {

	return false, nil
}
