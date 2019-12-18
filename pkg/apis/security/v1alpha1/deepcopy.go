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

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto for WhiteListItem
func (in *WhiteListItem) DeepCopyInto(out *WhiteListItem) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = WhiteListSpec{
		Creator: in.Spec.Creator,
		Reason:  in.Spec.Reason,
		Regexp:  in.Spec.Regexp,
	}
}

// DeepCopyObject for WhiteListItem
func (in *WhiteListItem) DeepCopyObject() runtime.Object {
	out := WhiteListItem{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject for WhiteList
func (in *WhiteListItemList) DeepCopyObject() runtime.Object {
	out := WhiteListItemList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]WhiteListItem, len(in.Items))

		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

// DeepCopyInto Audit
func (in *Audit) DeepCopyInto(out *Audit) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = AuditSpec{
		ReleaseName: in.Spec.ReleaseName,
		Resource:    in.Spec.Resource,
		Images:      in.Spec.Images,
		Result:      in.Spec.Result,
		Action:      in.Spec.Action,
	}
}

// DeepCopyObject for Audit
func (in *Audit) DeepCopyObject() runtime.Object {
	out := Audit{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject for AuditList
func (in *AuditList) DeepCopyObject() runtime.Object {
	out := AuditList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Audit, len(in.Items))

		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
