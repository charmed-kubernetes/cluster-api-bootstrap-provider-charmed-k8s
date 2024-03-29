//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfig) DeepCopyInto(out *CharmedK8sConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfig.
func (in *CharmedK8sConfig) DeepCopy() *CharmedK8sConfig {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CharmedK8sConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigList) DeepCopyInto(out *CharmedK8sConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CharmedK8sConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigList.
func (in *CharmedK8sConfigList) DeepCopy() *CharmedK8sConfigList {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CharmedK8sConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigSpec) DeepCopyInto(out *CharmedK8sConfigSpec) {
	*out = *in
	if in.ControlPlaneApplications != nil {
		in, out := &in.ControlPlaneApplications, &out.ControlPlaneApplications
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.WorkerApplications != nil {
		in, out := &in.WorkerApplications, &out.WorkerApplications
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigSpec.
func (in *CharmedK8sConfigSpec) DeepCopy() *CharmedK8sConfigSpec {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigStatus) DeepCopyInto(out *CharmedK8sConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigStatus.
func (in *CharmedK8sConfigStatus) DeepCopy() *CharmedK8sConfigStatus {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigTemplate) DeepCopyInto(out *CharmedK8sConfigTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigTemplate.
func (in *CharmedK8sConfigTemplate) DeepCopy() *CharmedK8sConfigTemplate {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CharmedK8sConfigTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigTemplateList) DeepCopyInto(out *CharmedK8sConfigTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CharmedK8sConfigTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigTemplateList.
func (in *CharmedK8sConfigTemplateList) DeepCopy() *CharmedK8sConfigTemplateList {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CharmedK8sConfigTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigTemplateResource) DeepCopyInto(out *CharmedK8sConfigTemplateResource) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigTemplateResource.
func (in *CharmedK8sConfigTemplateResource) DeepCopy() *CharmedK8sConfigTemplateResource {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CharmedK8sConfigTemplateSpec) DeepCopyInto(out *CharmedK8sConfigTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CharmedK8sConfigTemplateSpec.
func (in *CharmedK8sConfigTemplateSpec) DeepCopy() *CharmedK8sConfigTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(CharmedK8sConfigTemplateSpec)
	in.DeepCopyInto(out)
	return out
}
