//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v3 "github.com/Masterminds/semver/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AvailableModule) DeepCopyInto(out *AvailableModule) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AvailableModule.
func (in *AvailableModule) DeepCopy() *AvailableModule {
	if in == nil {
		return nil
	}
	out := new(AvailableModule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Duration) DeepCopyInto(out *Duration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Duration.
func (in *Duration) DeepCopy() *Duration {
	if in == nil {
		return nil
	}
	out := new(Duration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Module) DeepCopyInto(out *Module) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Properties = in.Properties
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Module.
func (in *Module) DeepCopy() *Module {
	if in == nil {
		return nil
	}
	out := new(Module)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Module) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleConfig) DeepCopyInto(out *ModuleConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleConfig.
func (in *ModuleConfig) DeepCopy() *ModuleConfig {
	if in == nil {
		return nil
	}
	out := new(ModuleConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleConfigList) DeepCopyInto(out *ModuleConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModuleConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleConfigList.
func (in *ModuleConfigList) DeepCopy() *ModuleConfigList {
	if in == nil {
		return nil
	}
	out := new(ModuleConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleConfigSpec) DeepCopyInto(out *ModuleConfigSpec) {
	*out = *in
	in.Settings.DeepCopyInto(&out.Settings)
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleConfigSpec.
func (in *ModuleConfigSpec) DeepCopy() *ModuleConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ModuleConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleConfigStatus) DeepCopyInto(out *ModuleConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleConfigStatus.
func (in *ModuleConfigStatus) DeepCopy() *ModuleConfigStatus {
	if in == nil {
		return nil
	}
	out := new(ModuleConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleError) DeepCopyInto(out *ModuleError) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleError.
func (in *ModuleError) DeepCopy() *ModuleError {
	if in == nil {
		return nil
	}
	out := new(ModuleError)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleList) DeepCopyInto(out *ModuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Module, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleList.
func (in *ModuleList) DeepCopy() *ModuleList {
	if in == nil {
		return nil
	}
	out := new(ModuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleProperties) DeepCopyInto(out *ModuleProperties) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleProperties.
func (in *ModuleProperties) DeepCopy() *ModuleProperties {
	if in == nil {
		return nil
	}
	out := new(ModuleProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModulePullOverride) DeepCopyInto(out *ModulePullOverride) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModulePullOverride.
func (in *ModulePullOverride) DeepCopy() *ModulePullOverride {
	if in == nil {
		return nil
	}
	out := new(ModulePullOverride)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModulePullOverride) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModulePullOverrideKind) DeepCopyInto(out *ModulePullOverrideKind) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModulePullOverrideKind.
func (in *ModulePullOverrideKind) DeepCopy() *ModulePullOverrideKind {
	if in == nil {
		return nil
	}
	out := new(ModulePullOverrideKind)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModulePullOverrideList) DeepCopyInto(out *ModulePullOverrideList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModulePullOverride, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModulePullOverrideList.
func (in *ModulePullOverrideList) DeepCopy() *ModulePullOverrideList {
	if in == nil {
		return nil
	}
	out := new(ModulePullOverrideList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModulePullOverrideList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModulePullOverrideSpec) DeepCopyInto(out *ModulePullOverrideSpec) {
	*out = *in
	out.ScanInterval = in.ScanInterval
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModulePullOverrideSpec.
func (in *ModulePullOverrideSpec) DeepCopy() *ModulePullOverrideSpec {
	if in == nil {
		return nil
	}
	out := new(ModulePullOverrideSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModulePullOverrideStatus) DeepCopyInto(out *ModulePullOverrideStatus) {
	*out = *in
	in.RenewAt.DeepCopyInto(&out.RenewAt)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModulePullOverrideStatus.
func (in *ModulePullOverrideStatus) DeepCopy() *ModulePullOverrideStatus {
	if in == nil {
		return nil
	}
	out := new(ModulePullOverrideStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleRelease) DeepCopyInto(out *ModuleRelease) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleRelease.
func (in *ModuleRelease) DeepCopy() *ModuleRelease {
	if in == nil {
		return nil
	}
	out := new(ModuleRelease)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleRelease) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleReleaseList) DeepCopyInto(out *ModuleReleaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModuleRelease, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleReleaseList.
func (in *ModuleReleaseList) DeepCopy() *ModuleReleaseList {
	if in == nil {
		return nil
	}
	out := new(ModuleReleaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleReleaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleReleaseSpec) DeepCopyInto(out *ModuleReleaseSpec) {
	*out = *in
	if in.Version != nil {
		in, out := &in.Version, &out.Version
		*out = new(v3.Version)
		**out = **in
	}
	if in.ApplyAfter != nil {
		in, out := &in.ApplyAfter, &out.ApplyAfter
		*out = (*in).DeepCopy()
	}
	if in.Requirements != nil {
		in, out := &in.Requirements, &out.Requirements
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleReleaseSpec.
func (in *ModuleReleaseSpec) DeepCopy() *ModuleReleaseSpec {
	if in == nil {
		return nil
	}
	out := new(ModuleReleaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleReleaseStatus) DeepCopyInto(out *ModuleReleaseStatus) {
	*out = *in
	in.TransitionTime.DeepCopyInto(&out.TransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleReleaseStatus.
func (in *ModuleReleaseStatus) DeepCopy() *ModuleReleaseStatus {
	if in == nil {
		return nil
	}
	out := new(ModuleReleaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleSource) DeepCopyInto(out *ModuleSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleSource.
func (in *ModuleSource) DeepCopy() *ModuleSource {
	if in == nil {
		return nil
	}
	out := new(ModuleSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleSourceList) DeepCopyInto(out *ModuleSourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModuleSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleSourceList.
func (in *ModuleSourceList) DeepCopy() *ModuleSourceList {
	if in == nil {
		return nil
	}
	out := new(ModuleSourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleSourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleSourceSpec) DeepCopyInto(out *ModuleSourceSpec) {
	*out = *in
	out.Registry = in.Registry
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleSourceSpec.
func (in *ModuleSourceSpec) DeepCopy() *ModuleSourceSpec {
	if in == nil {
		return nil
	}
	out := new(ModuleSourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleSourceSpecRegistry) DeepCopyInto(out *ModuleSourceSpecRegistry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleSourceSpecRegistry.
func (in *ModuleSourceSpecRegistry) DeepCopy() *ModuleSourceSpecRegistry {
	if in == nil {
		return nil
	}
	out := new(ModuleSourceSpecRegistry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleSourceStatus) DeepCopyInto(out *ModuleSourceStatus) {
	*out = *in
	in.SyncTime.DeepCopyInto(&out.SyncTime)
	if in.AvailableModules != nil {
		in, out := &in.AvailableModules, &out.AvailableModules
		*out = make([]AvailableModule, len(*in))
		copy(*out, *in)
	}
	if in.ModuleErrors != nil {
		in, out := &in.ModuleErrors, &out.ModuleErrors
		*out = make([]ModuleError, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleSourceStatus.
func (in *ModuleSourceStatus) DeepCopy() *ModuleSourceStatus {
	if in == nil {
		return nil
	}
	out := new(ModuleSourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleUpdatePolicy) DeepCopyInto(out *ModuleUpdatePolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleUpdatePolicy.
func (in *ModuleUpdatePolicy) DeepCopy() *ModuleUpdatePolicy {
	if in == nil {
		return nil
	}
	out := new(ModuleUpdatePolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleUpdatePolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleUpdatePolicyList) DeepCopyInto(out *ModuleUpdatePolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModuleUpdatePolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleUpdatePolicyList.
func (in *ModuleUpdatePolicyList) DeepCopy() *ModuleUpdatePolicyList {
	if in == nil {
		return nil
	}
	out := new(ModuleUpdatePolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModuleUpdatePolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleUpdatePolicySpec) DeepCopyInto(out *ModuleUpdatePolicySpec) {
	*out = *in
	in.Update.DeepCopyInto(&out.Update)
	in.ModuleReleaseSelector.DeepCopyInto(&out.ModuleReleaseSelector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleUpdatePolicySpec.
func (in *ModuleUpdatePolicySpec) DeepCopy() *ModuleUpdatePolicySpec {
	if in == nil {
		return nil
	}
	out := new(ModuleUpdatePolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleUpdatePolicySpecReleaseSelector) DeepCopyInto(out *ModuleUpdatePolicySpecReleaseSelector) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleUpdatePolicySpecReleaseSelector.
func (in *ModuleUpdatePolicySpecReleaseSelector) DeepCopy() *ModuleUpdatePolicySpecReleaseSelector {
	if in == nil {
		return nil
	}
	out := new(ModuleUpdatePolicySpecReleaseSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModuleUpdatePolicySpecUpdate) DeepCopyInto(out *ModuleUpdatePolicySpecUpdate) {
	*out = *in
	out.Windows = in.Windows.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModuleUpdatePolicySpecUpdate.
func (in *ModuleUpdatePolicySpecUpdate) DeepCopy() *ModuleUpdatePolicySpecUpdate {
	if in == nil {
		return nil
	}
	out := new(ModuleUpdatePolicySpecUpdate)
	in.DeepCopyInto(out)
	return out
}
