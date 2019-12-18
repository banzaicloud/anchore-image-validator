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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName for crd
const GroupName = "security.banzaicloud.com"

// GroupVersion for crd
const GroupVersion = "v1alpha1"

// nolint: gochecknoglobals
var (
	// SchemeGroupVersion for crd
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	// SchemeBuilder for crd
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme for crd
	AddToScheme = SchemeBuilder.AddToScheme
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
