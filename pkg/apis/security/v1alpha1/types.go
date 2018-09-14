package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type WhiteListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []WhiteList `json:"items"`
}

type WhiteList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec WhiteListSpec `json:"spec"`
}

type WhiteListSpec struct {
	ReleaseName string `json:"releaseName"`
	Creator     string `json:"creator"`
	Reason      string `json:"reason"`
}
