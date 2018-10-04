package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName for crd
const GroupName = "security.banzaicloud.com"
// GroupVersion for crd
const GroupVersion = "v1alpha1"

// SchemeGroupVersion for crd
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	// SchemeBuilder for crd
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme for crd
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&WhiteListItem{},
		&WhiteListItemList{},
		&AuditList{},
		&Audit{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
