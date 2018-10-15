package main

import (
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
	image       []string
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
	if f {
		logrus.WithFields(logrus.Fields{
			"FakeRelease": true,
		}).Info("Missing release label, using PodName")
		for _, res := range wl {
			fakeRelease := string(res.ObjectMeta.Name + "-")
			if strings.Contains(r, fakeRelease) {
				return true
			}
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"release": r,
		}).Info("Check whitelist")
		for _, res := range wl {
			if r == res.ObjectMeta.Name {
				return true
			}
		}
	}
	return false
}

func createAudit(a auditInfo) {
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
			Image:       a.image,
			Result:      a.result,
			Action:      a.action,
		},
		Status: v1alpha1.AuditStatus{
			State: a.state,
		},
	}
	auditCR.SetOwnerReferences(a.owners)
	audit, err := securityClientSet.Audits("default").Create(auditCR)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.WithFields(logrus.Fields{
			"Audit": audit,
		}).Debug("Created Audits")
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
