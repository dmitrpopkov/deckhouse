//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021 Flant JSC

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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BootstrapStatus) DeepCopyInto(out *BootstrapStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BootstrapStatus.
func (in *BootstrapStatus) DeepCopy() *BootstrapStatus {
	if in == nil {
		return nil
	}
	out := new(BootstrapStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CurrentStatus) DeepCopyInto(out *CurrentStatus) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CurrentStatus.
func (in *CurrentStatus) DeepCopy() *CurrentStatus {
	if in == nil {
		return nil
	}
	out := new(CurrentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceClaim) DeepCopyInto(out *InstanceClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceClaim.
func (in *InstanceClaim) DeepCopy() *InstanceClaim {
	if in == nil {
		return nil
	}
	out := new(InstanceClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceClaimOperation) DeepCopyInto(out *InstanceClaimOperation) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceClaimOperation.
func (in *InstanceClaimOperation) DeepCopy() *InstanceClaimOperation {
	if in == nil {
		return nil
	}
	out := new(InstanceClaimOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceClaimStatus) DeepCopyInto(out *InstanceClaimStatus) {
	*out = *in
	out.NodeRef = in.NodeRef
	out.MachineRef = in.MachineRef
	in.CurrentStatus.DeepCopyInto(&out.CurrentStatus)
	in.LastOperation.DeepCopyInto(&out.LastOperation)
	out.BootstrapStatus = in.BootstrapStatus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceClaimStatus.
func (in *InstanceClaimStatus) DeepCopy() *InstanceClaimStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LastOperation) DeepCopyInto(out *LastOperation) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LastOperation.
func (in *LastOperation) DeepCopy() *LastOperation {
	if in == nil {
		return nil
	}
	out := new(LastOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineRef) DeepCopyInto(out *MachineRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineRef.
func (in *MachineRef) DeepCopy() *MachineRef {
	if in == nil {
		return nil
	}
	out := new(MachineRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeRef) DeepCopyInto(out *NodeRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeRef.
func (in *NodeRef) DeepCopy() *NodeRef {
	if in == nil {
		return nil
	}
	out := new(NodeRef)
	in.DeepCopyInto(out)
	return out
}
