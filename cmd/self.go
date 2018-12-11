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
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/goph/emperror"
	"github.com/sirupsen/logrus"
	admissionV1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionClient "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	clientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func createValidatingWebhook(c *clientV1.CoreV1Client) *admissionV1beta1.ValidatingWebhookConfiguration {

	path := path.Join("/apis", apiServiceGroup, apiServiceVersion, apiServiceResource)
	webHookName := fmt.Sprintf("%s.%s", anchoreReleaseName, apiServiceGroup)
	ownerref, caBundle, err := getSelf(c)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	rule := admissionV1beta1.Rule{
		APIGroups:   []string{""},
		APIVersions: []string{"*"},
		Resources:   []string{"pods"},
	}

	rulesWithOperations := admissionV1beta1.RuleWithOperations{
		Operations: []admissionV1beta1.OperationType{admissionV1beta1.Create},
		Rule:       rule,
	}

	failurePolicy := admissionV1beta1.Fail

	expression := metav1.LabelSelectorRequirement{
		Key:      "scan",
		Operator: metav1.LabelSelectorOpNotIn,
		Values:   []string{"noscan"},
	}

	nameSpaceSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{expression},
	}

	validatingWebhook := admissionV1beta1.Webhook{
		Name: webHookName,
		ClientConfig: admissionV1beta1.WebhookClientConfig{
			Service: &admissionV1beta1.ServiceReference{
				Namespace: "default",
				Name:      "kubernetes",
				Path:      &path,
			},
			CABundle: caBundle,
		},
		Rules:             []admissionV1beta1.RuleWithOperations{rulesWithOperations},
		FailurePolicy:     &failurePolicy,
		NamespaceSelector: nameSpaceSelector,
	}

	validatingWebhookConfig := &admissionV1beta1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ValidatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: webHookName,
		},
		Webhooks: []admissionV1beta1.Webhook{validatingWebhook},
	}

	validatingWebhookConfig.SetOwnerReferences(ownerref)

	return validatingWebhookConfig
}

func installValidatingWebhookConfig(c *rest.Config) error {
	coreClientSet, err := clientV1.NewForConfig(c)
	if err != nil {
		logrus.Error(err)
	}
	validatingWebhookConfig := createValidatingWebhook(coreClientSet)
	if validatingWebhookConfig == nil {
		return emperror.Wrap(err, "cannot create ValidatingkWebhooConfiguration")
	}
	admissionClientSet, err := admissionClient.NewForConfig(c)
	if err != nil {
		return emperror.Wrap(err, "cannot create admission registration client")
	}
	validatingInt := admissionClientSet.ValidatingWebhookConfigurations()
	_, err = validatingInt.Create(validatingWebhookConfig)
	if err != nil {
		return emperror.Wrap(err, "cannot install ValidatingWebhookConfiguration")
	}
	return nil
}

func getSelf(c *clientV1.CoreV1Client) ([]metav1.OwnerReference, []byte, error) {
	podName, _ := os.Hostname()
	if kubernetesNameSpace == "" {
		return nil, nil, errors.New("not defined KUBERNETES_NAMESPACE env")
	}
	podDetail, err := c.Pods(kubernetesNameSpace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, emperror.Wrap(err, "unable to get self details")
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

	secretDetail, err := c.Secrets(kubernetesNameSpace).Get(anchoreReleaseName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, emperror.Wrap(err, "unable to get secretDetail")
	}
	caBundle := []byte(base64.StdEncoding.EncodeToString(secretDetail.Data["caCert"]))

	return []metav1.OwnerReference{owner}, caBundle, nil
}
