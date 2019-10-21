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

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"emperror.dev/errors"
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func createValidatingWebhook(c client.Client) (*admissionv1beta1.ValidatingWebhookConfiguration, error) {

	path := path.Join("/apis", apiServiceGroup, apiServiceVersion, apiServiceResource)
	webHookName := fmt.Sprintf("%s.%s", anchoreReleaseName, apiServiceGroup)
	ownerref, caBundle, err := getSelf(c)
	if err != nil {
		return nil, errors.WrapIf(err, "unable to get self object")
	}
	rule := admissionv1beta1.Rule{
		APIGroups:   []string{""},
		APIVersions: []string{"*"},
		Resources:   []string{"pods"},
	}

	rulesWithOperations := admissionv1beta1.RuleWithOperations{
		Operations: []admissionv1beta1.OperationType{admissionv1beta1.Create},
		Rule:       rule,
	}

	failurePolicy := admissionv1beta1.Fail

	selectorOperator := metav1.LabelSelectorOpNotIn
	selectorValues := []string{"noscan"}

	if namespaceSelector == "include" {
		selectorOperator = metav1.LabelSelectorOpIn
		selectorValues = []string{"scan"}
	}

	expression := metav1.LabelSelectorRequirement{
		Key:      "scan",
		Operator: selectorOperator,
		Values:   selectorValues,
	}

	nameSpaceSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{expression},
	}

	validatingWebhook := admissionv1beta1.ValidatingWebhook{
		Name: webHookName,
		ClientConfig: admissionv1beta1.WebhookClientConfig{
			Service: &admissionv1beta1.ServiceReference{
				Namespace: "default",
				Name:      "kubernetes",
				Path:      &path,
			},
			CABundle: caBundle,
		},
		Rules:             []admissionv1beta1.RuleWithOperations{rulesWithOperations},
		FailurePolicy:     &failurePolicy,
		NamespaceSelector: nameSpaceSelector,
	}

	validatingWebhookConfig := &admissionv1beta1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ValidatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: webHookName,
		},
		Webhooks: []admissionv1beta1.ValidatingWebhook{validatingWebhook},
	}

	validatingWebhookConfig.SetOwnerReferences(ownerref)

	return validatingWebhookConfig, nil
}

func installValidatingWebhookConfig(c client.Client) error {
	validatingWebhookConfig, err := createValidatingWebhook(c)
	if err != nil {
		return errors.WrapIf(err, "cannot create ValidatingkWebhooConfiguration")
	}

	err = c.Create(context.Background(), validatingWebhookConfig)
	if err != nil {
		return errors.WrapIf(err, "cannot install ValidatingWebhookConfiguration")
	}
	return nil
}

func getSelf(c client.Client) ([]metav1.OwnerReference, []byte, error) {
	podName, _ := os.Hostname()
	if kubernetesNameSpace == "" {
		return nil, nil, errors.New("not defined KUBERNETES_NAMESPACE env")
	}
	podDetail := &corev1.Pod{}
	err := c.Get(context.Background(), client.ObjectKey{
		Namespace: kubernetesNameSpace,
		Name:      podName,
	}, podDetail)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "unable to get self details")
	}

	if anchoreReleaseName == "" {
		return nil, nil, errors.New("not defined ANCHORE_RELEASE_NAME env")
	}

	owner := metav1.OwnerReference{
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       podName,
		UID:        podDetail.ObjectMeta.UID,
	}

	secretDetail := &corev1.Secret{}
	err = c.Get(context.Background(), client.ObjectKey{
		Namespace: kubernetesNameSpace,
		Name:      anchoreReleaseName,
	}, secretDetail)
	if err != nil {
		return nil, nil, errors.WrapIf(err, "unable to get secretDetail")
	}
	caBundle := []byte(base64.StdEncoding.EncodeToString(secretDetail.Data["caCert"]))

	return []metav1.OwnerReference{owner}, caBundle, nil
}
