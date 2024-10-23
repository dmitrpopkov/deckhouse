// Copyright 2024 Flant JSC
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: abort.proto

package dhctl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AbortRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*AbortRequest_Start
	//	*AbortRequest_Continue
	Message isAbortRequest_Message `protobuf_oneof:"message"`
}

func (x *AbortRequest) Reset() {
	*x = AbortRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortRequest) ProtoMessage() {}

func (x *AbortRequest) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortRequest.ProtoReflect.Descriptor instead.
func (*AbortRequest) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{0}
}

func (m *AbortRequest) GetMessage() isAbortRequest_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *AbortRequest) GetStart() *AbortStart {
	if x, ok := x.GetMessage().(*AbortRequest_Start); ok {
		return x.Start
	}
	return nil
}

func (x *AbortRequest) GetContinue() *AbortContinue {
	if x, ok := x.GetMessage().(*AbortRequest_Continue); ok {
		return x.Continue
	}
	return nil
}

type isAbortRequest_Message interface {
	isAbortRequest_Message()
}

type AbortRequest_Start struct {
	Start *AbortStart `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type AbortRequest_Continue struct {
	Continue *AbortContinue `protobuf:"bytes,2,opt,name=continue,proto3,oneof"`
}

func (*AbortRequest_Start) isAbortRequest_Message() {}

func (*AbortRequest_Continue) isAbortRequest_Message() {}

type AbortResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*AbortResponse_Result
	//	*AbortResponse_PhaseEnd
	//	*AbortResponse_Logs
	Message isAbortResponse_Message `protobuf_oneof:"message"`
}

func (x *AbortResponse) Reset() {
	*x = AbortResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortResponse) ProtoMessage() {}

func (x *AbortResponse) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortResponse.ProtoReflect.Descriptor instead.
func (*AbortResponse) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{1}
}

func (m *AbortResponse) GetMessage() isAbortResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *AbortResponse) GetResult() *AbortResult {
	if x, ok := x.GetMessage().(*AbortResponse_Result); ok {
		return x.Result
	}
	return nil
}

func (x *AbortResponse) GetPhaseEnd() *AbortPhaseEnd {
	if x, ok := x.GetMessage().(*AbortResponse_PhaseEnd); ok {
		return x.PhaseEnd
	}
	return nil
}

func (x *AbortResponse) GetLogs() *Logs {
	if x, ok := x.GetMessage().(*AbortResponse_Logs); ok {
		return x.Logs
	}
	return nil
}

type isAbortResponse_Message interface {
	isAbortResponse_Message()
}

type AbortResponse_Result struct {
	Result *AbortResult `protobuf:"bytes,1,opt,name=result,proto3,oneof"`
}

type AbortResponse_PhaseEnd struct {
	PhaseEnd *AbortPhaseEnd `protobuf:"bytes,2,opt,name=phase_end,json=phaseEnd,proto3,oneof"`
}

type AbortResponse_Logs struct {
	Logs *Logs `protobuf:"bytes,3,opt,name=logs,proto3,oneof"`
}

func (*AbortResponse_Result) isAbortResponse_Message() {}

func (*AbortResponse_PhaseEnd) isAbortResponse_Message() {}

func (*AbortResponse_Logs) isAbortResponse_Message() {}

type AbortStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConnectionConfig              string             `protobuf:"bytes,1,opt,name=connection_config,json=connectionConfig,proto3" json:"connection_config,omitempty"`
	InitConfig                    string             `protobuf:"bytes,2,opt,name=init_config,json=initConfig,proto3" json:"init_config,omitempty"`
	ClusterConfig                 string             `protobuf:"bytes,3,opt,name=cluster_config,json=clusterConfig,proto3" json:"cluster_config,omitempty"`
	ProviderSpecificClusterConfig string             `protobuf:"bytes,4,opt,name=provider_specific_cluster_config,json=providerSpecificClusterConfig,proto3" json:"provider_specific_cluster_config,omitempty"`
	InitResources                 string             `protobuf:"bytes,5,opt,name=init_resources,json=initResources,proto3" json:"init_resources,omitempty"`
	Resources                     string             `protobuf:"bytes,6,opt,name=resources,proto3" json:"resources,omitempty"`
	State                         string             `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	Options                       *AbortStartOptions `protobuf:"bytes,8,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *AbortStart) Reset() {
	*x = AbortStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortStart) ProtoMessage() {}

func (x *AbortStart) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortStart.ProtoReflect.Descriptor instead.
func (*AbortStart) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{2}
}

func (x *AbortStart) GetConnectionConfig() string {
	if x != nil {
		return x.ConnectionConfig
	}
	return ""
}

func (x *AbortStart) GetInitConfig() string {
	if x != nil {
		return x.InitConfig
	}
	return ""
}

func (x *AbortStart) GetClusterConfig() string {
	if x != nil {
		return x.ClusterConfig
	}
	return ""
}

func (x *AbortStart) GetProviderSpecificClusterConfig() string {
	if x != nil {
		return x.ProviderSpecificClusterConfig
	}
	return ""
}

func (x *AbortStart) GetInitResources() string {
	if x != nil {
		return x.InitResources
	}
	return ""
}

func (x *AbortStart) GetResources() string {
	if x != nil {
		return x.Resources
	}
	return ""
}

func (x *AbortStart) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *AbortStart) GetOptions() *AbortStartOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

type AbortPhaseEnd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CompletedPhase      string            `protobuf:"bytes,1,opt,name=completed_phase,json=completedPhase,proto3" json:"completed_phase,omitempty"`
	CompletedPhaseState map[string][]byte `protobuf:"bytes,2,rep,name=completed_phase_state,json=completedPhaseState,proto3" json:"completed_phase_state,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NextPhase           string            `protobuf:"bytes,3,opt,name=next_phase,json=nextPhase,proto3" json:"next_phase,omitempty"`
	NextPhaseCritical   bool              `protobuf:"varint,4,opt,name=next_phase_critical,json=nextPhaseCritical,proto3" json:"next_phase_critical,omitempty"`
}

func (x *AbortPhaseEnd) Reset() {
	*x = AbortPhaseEnd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortPhaseEnd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortPhaseEnd) ProtoMessage() {}

func (x *AbortPhaseEnd) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortPhaseEnd.ProtoReflect.Descriptor instead.
func (*AbortPhaseEnd) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{3}
}

func (x *AbortPhaseEnd) GetCompletedPhase() string {
	if x != nil {
		return x.CompletedPhase
	}
	return ""
}

func (x *AbortPhaseEnd) GetCompletedPhaseState() map[string][]byte {
	if x != nil {
		return x.CompletedPhaseState
	}
	return nil
}

func (x *AbortPhaseEnd) GetNextPhase() string {
	if x != nil {
		return x.NextPhase
	}
	return ""
}

func (x *AbortPhaseEnd) GetNextPhaseCritical() bool {
	if x != nil {
		return x.NextPhaseCritical
	}
	return false
}

type AbortContinue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Continue Continue `protobuf:"varint,1,opt,name=continue,proto3,enum=dhctl.Continue" json:"continue,omitempty"`
	Err      string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *AbortContinue) Reset() {
	*x = AbortContinue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortContinue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortContinue) ProtoMessage() {}

func (x *AbortContinue) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortContinue.ProtoReflect.Descriptor instead.
func (*AbortContinue) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{4}
}

func (x *AbortContinue) GetContinue() Continue {
	if x != nil {
		return x.Continue
	}
	return Continue_CONTINUE_UNSPECIFIED
}

func (x *AbortContinue) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type AbortStartOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommanderMode    bool                 `protobuf:"varint,1,opt,name=commander_mode,json=commanderMode,proto3" json:"commander_mode,omitempty"`
	CommanderUuid    string               `protobuf:"bytes,2,opt,name=commander_uuid,json=commanderUuid,proto3" json:"commander_uuid,omitempty"`
	LogWidth         int32                `protobuf:"varint,3,opt,name=log_width,json=logWidth,proto3" json:"log_width,omitempty"`
	ResourcesTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=resources_timeout,json=resourcesTimeout,proto3" json:"resources_timeout,omitempty"`
	DeckhouseTimeout *durationpb.Duration `protobuf:"bytes,5,opt,name=deckhouse_timeout,json=deckhouseTimeout,proto3" json:"deckhouse_timeout,omitempty"`
	CommonOptions    *OperationOptions    `protobuf:"bytes,10,opt,name=common_options,json=commonOptions,proto3" json:"common_options,omitempty"`
}

func (x *AbortStartOptions) Reset() {
	*x = AbortStartOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortStartOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortStartOptions) ProtoMessage() {}

func (x *AbortStartOptions) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortStartOptions.ProtoReflect.Descriptor instead.
func (*AbortStartOptions) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{5}
}

func (x *AbortStartOptions) GetCommanderMode() bool {
	if x != nil {
		return x.CommanderMode
	}
	return false
}

func (x *AbortStartOptions) GetCommanderUuid() string {
	if x != nil {
		return x.CommanderUuid
	}
	return ""
}

func (x *AbortStartOptions) GetLogWidth() int32 {
	if x != nil {
		return x.LogWidth
	}
	return 0
}

func (x *AbortStartOptions) GetResourcesTimeout() *durationpb.Duration {
	if x != nil {
		return x.ResourcesTimeout
	}
	return nil
}

func (x *AbortStartOptions) GetDeckhouseTimeout() *durationpb.Duration {
	if x != nil {
		return x.DeckhouseTimeout
	}
	return nil
}

func (x *AbortStartOptions) GetCommonOptions() *OperationOptions {
	if x != nil {
		return x.CommonOptions
	}
	return nil
}

type AbortResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	Err   string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *AbortResult) Reset() {
	*x = AbortResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_abort_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbortResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbortResult) ProtoMessage() {}

func (x *AbortResult) ProtoReflect() protoreflect.Message {
	mi := &file_abort_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbortResult.ProtoReflect.Descriptor instead.
func (*AbortResult) Descriptor() ([]byte, []int) {
	return file_abort_proto_rawDescGZIP(), []int{6}
}

func (x *AbortResult) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *AbortResult) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_abort_proto protoreflect.FileDescriptor

var file_abort_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x62, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x64,
	0x68, 0x63, 0x74, 0x6c, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x78, 0x0a, 0x0c, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x29, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x32, 0x0a,
	0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x43, 0x6f, 0x6e,
	0x74, 0x69, 0x6e, 0x75, 0x65, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75,
	0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xa0, 0x01, 0x0a,
	0x0d, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x48, 0x00, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x33, 0x0a, 0x09,
	0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x50, 0x68, 0x61,
	0x73, 0x65, 0x45, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x08, 0x70, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e,
	0x64, 0x12, 0x21, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x4c, 0x6f, 0x67, 0x73, 0x48, 0x00, 0x52, 0x04,
	0x6c, 0x6f, 0x67, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0xd9, 0x02, 0x0a, 0x0a, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2b,
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x0a, 0x0b, 0x69,
	0x6e, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x69, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a, 0x0e,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x47, 0x0a, 0x20, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f,
	0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1d, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a, 0x0e,
	0x69, 0x6e, 0x69, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c,
	0x2e, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xb2, 0x02, 0x0a, 0x0d,
	0x41, 0x62, 0x6f, 0x72, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x27, 0x0a,
	0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x61, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x65, 0x64, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x41, 0x62,
	0x6f, 0x72, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x13, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x50,
	0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78,
	0x74, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x6e, 0x65, 0x78, 0x74,
	0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x68, 0x61, 0x73, 0x65,
	0x43, 0x72, 0x69, 0x74, 0x69, 0x63, 0x61, 0x6c, 0x1a, 0x46, 0x0a, 0x18, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x50, 0x68, 0x61, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x4e, 0x0a, 0x0d, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75,
	0x65, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x64, 0x68, 0x63, 0x74, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x69, 0x6e, 0x75, 0x65, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72,
	0x22, 0xce, 0x02, 0x0a, 0x11, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72,
	0x55, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x77, 0x69, 0x64, 0x74,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x57, 0x69, 0x64, 0x74,
	0x68, 0x12, 0x46, 0x0a, 0x11, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x46, 0x0a, 0x11, 0x64, 0x65, 0x63,
	0x6b, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x10, 0x64, 0x65, 0x63, 0x6b, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x12, 0x3e, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x68, 0x63, 0x74,
	0x6c, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x22, 0x35, 0x0a, 0x0b, 0x41, 0x62, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x62, 0x2f, 0x64,
	0x68, 0x63, 0x74, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_abort_proto_rawDescOnce sync.Once
	file_abort_proto_rawDescData = file_abort_proto_rawDesc
)

func file_abort_proto_rawDescGZIP() []byte {
	file_abort_proto_rawDescOnce.Do(func() {
		file_abort_proto_rawDescData = protoimpl.X.CompressGZIP(file_abort_proto_rawDescData)
	})
	return file_abort_proto_rawDescData
}

var file_abort_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_abort_proto_goTypes = []interface{}{
	(*AbortRequest)(nil),        // 0: dhctl.AbortRequest
	(*AbortResponse)(nil),       // 1: dhctl.AbortResponse
	(*AbortStart)(nil),          // 2: dhctl.AbortStart
	(*AbortPhaseEnd)(nil),       // 3: dhctl.AbortPhaseEnd
	(*AbortContinue)(nil),       // 4: dhctl.AbortContinue
	(*AbortStartOptions)(nil),   // 5: dhctl.AbortStartOptions
	(*AbortResult)(nil),         // 6: dhctl.AbortResult
	nil,                         // 7: dhctl.AbortPhaseEnd.CompletedPhaseStateEntry
	(*Logs)(nil),                // 8: dhctl.Logs
	(Continue)(0),               // 9: dhctl.Continue
	(*durationpb.Duration)(nil), // 10: google.protobuf.Duration
	(*OperationOptions)(nil),    // 11: dhctl.OperationOptions
}
var file_abort_proto_depIdxs = []int32{
	2,  // 0: dhctl.AbortRequest.start:type_name -> dhctl.AbortStart
	4,  // 1: dhctl.AbortRequest.continue:type_name -> dhctl.AbortContinue
	6,  // 2: dhctl.AbortResponse.result:type_name -> dhctl.AbortResult
	3,  // 3: dhctl.AbortResponse.phase_end:type_name -> dhctl.AbortPhaseEnd
	8,  // 4: dhctl.AbortResponse.logs:type_name -> dhctl.Logs
	5,  // 5: dhctl.AbortStart.options:type_name -> dhctl.AbortStartOptions
	7,  // 6: dhctl.AbortPhaseEnd.completed_phase_state:type_name -> dhctl.AbortPhaseEnd.CompletedPhaseStateEntry
	9,  // 7: dhctl.AbortContinue.continue:type_name -> dhctl.Continue
	10, // 8: dhctl.AbortStartOptions.resources_timeout:type_name -> google.protobuf.Duration
	10, // 9: dhctl.AbortStartOptions.deckhouse_timeout:type_name -> google.protobuf.Duration
	11, // 10: dhctl.AbortStartOptions.common_options:type_name -> dhctl.OperationOptions
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_abort_proto_init() }
func file_abort_proto_init() {
	if File_abort_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_abort_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortStart); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortPhaseEnd); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortContinue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortStartOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_abort_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbortResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_abort_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*AbortRequest_Start)(nil),
		(*AbortRequest_Continue)(nil),
	}
	file_abort_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*AbortResponse_Result)(nil),
		(*AbortResponse_PhaseEnd)(nil),
		(*AbortResponse_Logs)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_abort_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_abort_proto_goTypes,
		DependencyIndexes: file_abort_proto_depIdxs,
		MessageInfos:      file_abort_proto_msgTypes,
	}.Build()
	File_abort_proto = out.File
	file_abort_proto_rawDesc = nil
	file_abort_proto_goTypes = nil
	file_abort_proto_depIdxs = nil
}
