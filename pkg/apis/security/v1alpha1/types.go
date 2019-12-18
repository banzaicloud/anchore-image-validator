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

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// WhiteListItemList for whitelisting crd
type WhiteListItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []WhiteListItem `json:"items"`
}

// WhiteListItem for whitelisting crd
type WhiteListItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec WhiteListSpec `json:"spec"`
}

// WhiteListSpec for WhiteListItem
type WhiteListSpec struct {
	Creator string `json:"creator"`
	Reason  string `json:"reason"`
	Regexp  string `json:"regexp,omitempty"`
}

// AuditList for Scan events
type AuditList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Audit `json:"items"`
}

// Audit for AuditList
type Audit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuditSpec   `json:"spec"`
	Status AuditStatus `json:"status,omitempty"`
}

// AuditSpec for Audit
type AuditSpec struct {
	ReleaseName string       `json:"releaseName"`
	Resource    string       `json:"resource"`
	Result      []string     `json:"result"`
	Action      string       `json:"action"`
	Images      []AuditImage `json:"image"`
}

// AuditImage for AuditSpec
type AuditImage struct {
	ImageName   string `json:"imageName"`
	ImageTag    string `json:"imageTag"`
	ImageDigest string `json:"imageDigest"`
	LastUpdated string `json:"lastUpdated"`
}

// AuditStatus for AuditSpec
type AuditStatus struct {
	State string `json:"state"`
}
