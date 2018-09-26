package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// WhiteList for whitelisting crd
type WhiteList struct {
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
	ReleaseName string `json:"releaseName"`
	Creator     string `json:"creator"`
	Reason      string `json:"reason"`
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
	ReleaseName string   `json:"releaseName"`
	Resource    string   `json:"resource"`
	Result      []string `json:"result"`
	Action      string   `json:"action"`
	Image       []string `json:"image"`
}

// AuditStatus for AuditSpec
type AuditStatus struct {
	State string `json:"state"`
}
