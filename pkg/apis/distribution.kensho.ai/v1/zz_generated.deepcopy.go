// +build !ignore_autogenerated

/*
This is the test project from Kensho
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KDeployment) DeepCopyInto(out *KDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KDeployment.
func (in *KDeployment) DeepCopy() *KDeployment {
	if in == nil {
		return nil
	}
	out := new(KDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KDeploymentList) DeepCopyInto(out *KDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KDeploymentList.
func (in *KDeploymentList) DeepCopy() *KDeploymentList {
	if in == nil {
		return nil
	}
	out := new(KDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KDeploymentSpec) DeepCopyInto(out *KDeploymentSpec) {
	*out = *in
	if in.TotalReplicas != nil {
		in, out := &in.TotalReplicas, &out.TotalReplicas
		*out = new(int32)
		**out = **in
	}
	in.DeploymentTemplate.DeepCopyInto(&out.DeploymentTemplate)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KDeploymentSpec.
func (in *KDeploymentSpec) DeepCopy() *KDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(KDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KDeploymentStatus) DeepCopyInto(out *KDeploymentStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KDeploymentStatus.
func (in *KDeploymentStatus) DeepCopy() *KDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(KDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}