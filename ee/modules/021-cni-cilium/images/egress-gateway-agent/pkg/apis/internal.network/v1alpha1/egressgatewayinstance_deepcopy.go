/*
Copyright 2024 Flant JSC

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

package v1alpha1

import (
	. "github.com/deckhouse/deckhouse/egress-gateway-agent/pkg/apis/common"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *SDNInternalEgressGatewayInstance) DeepCopyInto(out *SDNInternalEgressGatewayInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SDNInternalEgressGatewayInstance.
func (in *SDNInternalEgressGatewayInstance) DeepCopy() *SDNInternalEgressGatewayInstance {
	if in == nil {
		return nil
	}
	out := new(SDNInternalEgressGatewayInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SDNInternalEgressGatewayInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SDNInternalEgressGatewayInstanceList) DeepCopyInto(out *SDNInternalEgressGatewayInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SDNInternalEgressGatewayInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SDNInternalEgressGatewayInstanceList.
func (in *SDNInternalEgressGatewayInstanceList) DeepCopy() *SDNInternalEgressGatewayInstanceList {
	if in == nil {
		return nil
	}
	out := new(SDNInternalEgressGatewayInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SDNInternalEgressGatewayInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SDNInternalEgressGatewayInstanceSpec) DeepCopyInto(out *SDNInternalEgressGatewayInstanceSpec) {
	*out = *in
	out.SourceIP = in.SourceIP
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SDNInternalEgressGatewayInstanceSpec.
func (in *SDNInternalEgressGatewayInstanceSpec) DeepCopy() *SDNInternalEgressGatewayInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(SDNInternalEgressGatewayInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SDNInternalEgressGatewayInstanceStatus) DeepCopyInto(out *SDNInternalEgressGatewayInstanceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ExtendedCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SDNInternalEgressGatewayInstanceStatus.
func (in *SDNInternalEgressGatewayInstanceStatus) DeepCopy() *SDNInternalEgressGatewayInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(SDNInternalEgressGatewayInstanceStatus)
	in.DeepCopyInto(out)
	return out
}
