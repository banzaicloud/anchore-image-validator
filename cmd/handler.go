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
	"regexp"
	"strings"

	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type auditInfo struct {
	name        string
	labels      map[string]string
	releaseName string
	resource    string
	images      []v1alpha1.AuditImage
	result      []string
	action      string
	state       string
	owners      []metav1.OwnerReference
}

func getReleaseName(labels map[string]string, p string) (string, bool) {
	release := labels["release"]
	if release != "" {
		logrus.WithFields(logrus.Fields{
			"release": release,
		}).Info("Check whitelist")
		return release, false
	}
	logrus.WithFields(logrus.Fields{
		"PodName": p,
	}).Info("Missing release label, using PodName")
	return p, true
}

func checkWhiteList(wl []v1alpha1.WhiteListItem, r string, f bool) bool {
	for _, res := range wl {
		if f {
			logrus.WithFields(logrus.Fields{
				"FakeRelease": true,
			}).Info("Missing release label, using PodName")
			fakeRelease := string(res.ObjectMeta.Name + "-")
			if strings.Contains(r, fakeRelease) {
				return true
			}
		}
		if r == res.ObjectMeta.Name {
			return true
		}
		match := regexpWhiteList(res)
		if match != nil {
			if match.MatchString(r) {
				return true
			}
		}
	}

	return false
}

func regexpWhiteList(wl v1alpha1.WhiteListItem) *regexp.Regexp {
	if wl.Spec.Regexp != "" {
		match, err := regexp.Compile(wl.Spec.Regexp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("regexp compile error")
			return nil
		}
		return match
	}
	return nil
}

func createOrUpdateAudit(a auditInfo) {
	auditCR := &v1alpha1.Audit{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Audit",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   a.name,
			Labels: a.labels,
		},
		Spec: v1alpha1.AuditSpec{
			ReleaseName: a.releaseName,
			Resource:    a.resource,
			Images:      a.images,
			Result:      a.result,
			Action:      a.action,
		},
		Status: v1alpha1.AuditStatus{
			State: a.state,
		},
	}
	actionByte := []byte(`{"spec":`)
	aSpec, er := json.Marshal(auditCR.Spec)
	if er != nil {
		logrus.Error(er)
	}
	actionByte = append(actionByte, aSpec...)
	tail := []byte(`}`)
	actionByte = append(actionByte, tail...)
	auditCR.SetOwnerReferences(a.owners)
	audit, err := securityClientSet.Audits("default").Create(auditCR)
	if err != nil {
		logrus.Error(err)
		audit, err = securityClientSet.Audits("default").Update(a.name, actionByte)
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.WithFields(logrus.Fields{
				"Audit": audit,
			}).Debug("Update Audit")
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"Audit": audit,
		}).Debug("Created Audit")
	}
}

func listAudits() {
	audits, err := securityClientSet.Audits("default").List(metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.WithFields(logrus.Fields{
			"Audits": audits,
		}).Info("Listing Audits")
	}
}
