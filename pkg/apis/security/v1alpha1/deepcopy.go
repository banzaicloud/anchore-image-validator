package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto for WhiteListItem
func (in *WhiteListItem) DeepCopyInto(out *WhiteListItem) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = WhiteListSpec{
		ReleaseName: in.Spec.ReleaseName,
		Creator:     in.Spec.Creator,
		Reason:      in.Spec.Reason,
	}
}

// DeepCopyObject for WhiteListItem
func (in *WhiteListItem) DeepCopyObject() runtime.Object {
	out := WhiteListItem{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject for WhiteList
func (in *WhiteList) DeepCopyObject() runtime.Object {
	out := WhiteList{}
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
		Image:       in.Spec.Image,
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
