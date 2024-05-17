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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: dhctl.proto

package dhctl

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DHCTL_Check_FullMethodName     = "/dhctl.DHCTL/Check"
	DHCTL_Bootstrap_FullMethodName = "/dhctl.DHCTL/Bootstrap"
	DHCTL_Destroy_FullMethodName   = "/dhctl.DHCTL/Destroy"
	DHCTL_Abort_FullMethodName     = "/dhctl.DHCTL/Abort"
	DHCTL_Converge_FullMethodName  = "/dhctl.DHCTL/Converge"
	DHCTL_Import_FullMethodName    = "/dhctl.DHCTL/Import"
)

// DHCTLClient is the client API for DHCTL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DHCTLClient interface {
	Check(ctx context.Context, opts ...grpc.CallOption) (DHCTL_CheckClient, error)
	Bootstrap(ctx context.Context, opts ...grpc.CallOption) (DHCTL_BootstrapClient, error)
	Destroy(ctx context.Context, opts ...grpc.CallOption) (DHCTL_DestroyClient, error)
	Abort(ctx context.Context, opts ...grpc.CallOption) (DHCTL_AbortClient, error)
	Converge(ctx context.Context, opts ...grpc.CallOption) (DHCTL_ConvergeClient, error)
	Import(ctx context.Context, opts ...grpc.CallOption) (DHCTL_ImportClient, error)
}

type dHCTLClient struct {
	cc grpc.ClientConnInterface
}

func NewDHCTLClient(cc grpc.ClientConnInterface) DHCTLClient {
	return &dHCTLClient{cc}
}

func (c *dHCTLClient) Check(ctx context.Context, opts ...grpc.CallOption) (DHCTL_CheckClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[0], DHCTL_Check_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLCheckClient{stream}
	return x, nil
}

type DHCTL_CheckClient interface {
	Send(*CheckRequest) error
	Recv() (*CheckResponse, error)
	grpc.ClientStream
}

type dHCTLCheckClient struct {
	grpc.ClientStream
}

func (x *dHCTLCheckClient) Send(m *CheckRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLCheckClient) Recv() (*CheckResponse, error) {
	m := new(CheckResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHCTLClient) Bootstrap(ctx context.Context, opts ...grpc.CallOption) (DHCTL_BootstrapClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[1], DHCTL_Bootstrap_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLBootstrapClient{stream}
	return x, nil
}

type DHCTL_BootstrapClient interface {
	Send(*BootstrapRequest) error
	Recv() (*BootstrapResponse, error)
	grpc.ClientStream
}

type dHCTLBootstrapClient struct {
	grpc.ClientStream
}

func (x *dHCTLBootstrapClient) Send(m *BootstrapRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLBootstrapClient) Recv() (*BootstrapResponse, error) {
	m := new(BootstrapResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHCTLClient) Destroy(ctx context.Context, opts ...grpc.CallOption) (DHCTL_DestroyClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[2], DHCTL_Destroy_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLDestroyClient{stream}
	return x, nil
}

type DHCTL_DestroyClient interface {
	Send(*DestroyRequest) error
	Recv() (*DestroyResponse, error)
	grpc.ClientStream
}

type dHCTLDestroyClient struct {
	grpc.ClientStream
}

func (x *dHCTLDestroyClient) Send(m *DestroyRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLDestroyClient) Recv() (*DestroyResponse, error) {
	m := new(DestroyResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHCTLClient) Abort(ctx context.Context, opts ...grpc.CallOption) (DHCTL_AbortClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[3], DHCTL_Abort_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLAbortClient{stream}
	return x, nil
}

type DHCTL_AbortClient interface {
	Send(*AbortRequest) error
	Recv() (*AbortResponse, error)
	grpc.ClientStream
}

type dHCTLAbortClient struct {
	grpc.ClientStream
}

func (x *dHCTLAbortClient) Send(m *AbortRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLAbortClient) Recv() (*AbortResponse, error) {
	m := new(AbortResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHCTLClient) Converge(ctx context.Context, opts ...grpc.CallOption) (DHCTL_ConvergeClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[4], DHCTL_Converge_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLConvergeClient{stream}
	return x, nil
}

type DHCTL_ConvergeClient interface {
	Send(*ConvergeRequest) error
	Recv() (*ConvergeResponse, error)
	grpc.ClientStream
}

type dHCTLConvergeClient struct {
	grpc.ClientStream
}

func (x *dHCTLConvergeClient) Send(m *ConvergeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLConvergeClient) Recv() (*ConvergeResponse, error) {
	m := new(ConvergeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHCTLClient) Import(ctx context.Context, opts ...grpc.CallOption) (DHCTL_ImportClient, error) {
	stream, err := c.cc.NewStream(ctx, &DHCTL_ServiceDesc.Streams[5], DHCTL_Import_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dHCTLImportClient{stream}
	return x, nil
}

type DHCTL_ImportClient interface {
	Send(*ImportRequest) error
	Recv() (*ImportResponse, error)
	grpc.ClientStream
}

type dHCTLImportClient struct {
	grpc.ClientStream
}

func (x *dHCTLImportClient) Send(m *ImportRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHCTLImportClient) Recv() (*ImportResponse, error) {
	m := new(ImportResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DHCTLServer is the server API for DHCTL service.
// All implementations must embed UnimplementedDHCTLServer
// for forward compatibility
type DHCTLServer interface {
	Check(DHCTL_CheckServer) error
	Bootstrap(DHCTL_BootstrapServer) error
	Destroy(DHCTL_DestroyServer) error
	Abort(DHCTL_AbortServer) error
	Converge(DHCTL_ConvergeServer) error
	Import(DHCTL_ImportServer) error
	mustEmbedUnimplementedDHCTLServer()
}

// UnimplementedDHCTLServer must be embedded to have forward compatible implementations.
type UnimplementedDHCTLServer struct {
}

func (UnimplementedDHCTLServer) Check(DHCTL_CheckServer) error {
	return status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedDHCTLServer) Bootstrap(DHCTL_BootstrapServer) error {
	return status.Errorf(codes.Unimplemented, "method Bootstrap not implemented")
}
func (UnimplementedDHCTLServer) Destroy(DHCTL_DestroyServer) error {
	return status.Errorf(codes.Unimplemented, "method Destroy not implemented")
}
func (UnimplementedDHCTLServer) Abort(DHCTL_AbortServer) error {
	return status.Errorf(codes.Unimplemented, "method Abort not implemented")
}
func (UnimplementedDHCTLServer) Converge(DHCTL_ConvergeServer) error {
	return status.Errorf(codes.Unimplemented, "method Converge not implemented")
}
func (UnimplementedDHCTLServer) Import(DHCTL_ImportServer) error {
	return status.Errorf(codes.Unimplemented, "method Import not implemented")
}
func (UnimplementedDHCTLServer) mustEmbedUnimplementedDHCTLServer() {}

// UnsafeDHCTLServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DHCTLServer will
// result in compilation errors.
type UnsafeDHCTLServer interface {
	mustEmbedUnimplementedDHCTLServer()
}

func RegisterDHCTLServer(s grpc.ServiceRegistrar, srv DHCTLServer) {
	s.RegisterService(&DHCTL_ServiceDesc, srv)
}

func _DHCTL_Check_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Check(&dHCTLCheckServer{stream})
}

type DHCTL_CheckServer interface {
	Send(*CheckResponse) error
	Recv() (*CheckRequest, error)
	grpc.ServerStream
}

type dHCTLCheckServer struct {
	grpc.ServerStream
}

func (x *dHCTLCheckServer) Send(m *CheckResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLCheckServer) Recv() (*CheckRequest, error) {
	m := new(CheckRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHCTL_Bootstrap_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Bootstrap(&dHCTLBootstrapServer{stream})
}

type DHCTL_BootstrapServer interface {
	Send(*BootstrapResponse) error
	Recv() (*BootstrapRequest, error)
	grpc.ServerStream
}

type dHCTLBootstrapServer struct {
	grpc.ServerStream
}

func (x *dHCTLBootstrapServer) Send(m *BootstrapResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLBootstrapServer) Recv() (*BootstrapRequest, error) {
	m := new(BootstrapRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHCTL_Destroy_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Destroy(&dHCTLDestroyServer{stream})
}

type DHCTL_DestroyServer interface {
	Send(*DestroyResponse) error
	Recv() (*DestroyRequest, error)
	grpc.ServerStream
}

type dHCTLDestroyServer struct {
	grpc.ServerStream
}

func (x *dHCTLDestroyServer) Send(m *DestroyResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLDestroyServer) Recv() (*DestroyRequest, error) {
	m := new(DestroyRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHCTL_Abort_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Abort(&dHCTLAbortServer{stream})
}

type DHCTL_AbortServer interface {
	Send(*AbortResponse) error
	Recv() (*AbortRequest, error)
	grpc.ServerStream
}

type dHCTLAbortServer struct {
	grpc.ServerStream
}

func (x *dHCTLAbortServer) Send(m *AbortResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLAbortServer) Recv() (*AbortRequest, error) {
	m := new(AbortRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHCTL_Converge_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Converge(&dHCTLConvergeServer{stream})
}

type DHCTL_ConvergeServer interface {
	Send(*ConvergeResponse) error
	Recv() (*ConvergeRequest, error)
	grpc.ServerStream
}

type dHCTLConvergeServer struct {
	grpc.ServerStream
}

func (x *dHCTLConvergeServer) Send(m *ConvergeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLConvergeServer) Recv() (*ConvergeRequest, error) {
	m := new(ConvergeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHCTL_Import_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHCTLServer).Import(&dHCTLImportServer{stream})
}

type DHCTL_ImportServer interface {
	Send(*ImportResponse) error
	Recv() (*ImportRequest, error)
	grpc.ServerStream
}

type dHCTLImportServer struct {
	grpc.ServerStream
}

func (x *dHCTLImportServer) Send(m *ImportResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHCTLImportServer) Recv() (*ImportRequest, error) {
	m := new(ImportRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DHCTL_ServiceDesc is the grpc.ServiceDesc for DHCTL service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DHCTL_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dhctl.DHCTL",
	HandlerType: (*DHCTLServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Check",
			Handler:       _DHCTL_Check_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Bootstrap",
			Handler:       _DHCTL_Bootstrap_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Destroy",
			Handler:       _DHCTL_Destroy_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Abort",
			Handler:       _DHCTL_Abort_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Converge",
			Handler:       _DHCTL_Converge_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Import",
			Handler:       _DHCTL_Import_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "dhctl.proto",
}
