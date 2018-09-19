package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func (in *WhiteList) DeepCopyInto(out *WhiteList) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = WhiteListSpec{
		ReleaseName: in.Spec.ReleaseName,
		Creator:     in.Spec.Creator,
		Reason:      in.Spec.Reason,
	}
}

func (in *WhiteList) DeepCopyObject() runtime.Object {
	out := WhiteList{}
	in.DeepCopyInto(&out)

	return &out
}

func (in *WhiteListList) DeepCopyObject() runtime.Object {
	out := WhiteListList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]WhiteList, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

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

func (in *Audit) DeepCopyObject() runtime.Object {
	out := Audit{}
	in.DeepCopyInto(&out)

	return &out
}

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
