// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: pikvm_automator.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	PiKvmAutomator_CommandList_FullMethodName   = "/pikvm_automator.PiKvmAutomator/CommandList"
	PiKvmAutomator_CallCommand_FullMethodName   = "/pikvm_automator.PiKvmAutomator/CallCommand"
	PiKvmAutomator_DeleteCommand_FullMethodName = "/pikvm_automator.PiKvmAutomator/DeleteCommand"
	PiKvmAutomator_CreateCommand_FullMethodName = "/pikvm_automator.PiKvmAutomator/CreateCommand"
)

// PiKvmAutomatorClient is the client API for PiKvmAutomator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PiKvmAutomatorClient interface {
	CommandList(ctx context.Context, in *CommandListRequest, opts ...grpc.CallOption) (*CommandListResponse, error)
	CallCommand(ctx context.Context, in *CallCommandRequest, opts ...grpc.CallOption) (*CallCommandResponse, error)
	DeleteCommand(ctx context.Context, in *DeleteCommandRequest, opts ...grpc.CallOption) (*DeleteCommandResponse, error)
	CreateCommand(ctx context.Context, in *CreateCommandRequest, opts ...grpc.CallOption) (*CreateCommandResponse, error)
}

type piKvmAutomatorClient struct {
	cc grpc.ClientConnInterface
}

func NewPiKvmAutomatorClient(cc grpc.ClientConnInterface) PiKvmAutomatorClient {
	return &piKvmAutomatorClient{cc}
}

func (c *piKvmAutomatorClient) CommandList(ctx context.Context, in *CommandListRequest, opts ...grpc.CallOption) (*CommandListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandListResponse)
	err := c.cc.Invoke(ctx, PiKvmAutomator_CommandList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piKvmAutomatorClient) CallCommand(ctx context.Context, in *CallCommandRequest, opts ...grpc.CallOption) (*CallCommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CallCommandResponse)
	err := c.cc.Invoke(ctx, PiKvmAutomator_CallCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piKvmAutomatorClient) DeleteCommand(ctx context.Context, in *DeleteCommandRequest, opts ...grpc.CallOption) (*DeleteCommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCommandResponse)
	err := c.cc.Invoke(ctx, PiKvmAutomator_DeleteCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *piKvmAutomatorClient) CreateCommand(ctx context.Context, in *CreateCommandRequest, opts ...grpc.CallOption) (*CreateCommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCommandResponse)
	err := c.cc.Invoke(ctx, PiKvmAutomator_CreateCommand_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PiKvmAutomatorServer is the server API for PiKvmAutomator service.
// All implementations must embed UnimplementedPiKvmAutomatorServer
// for forward compatibility.
type PiKvmAutomatorServer interface {
	CommandList(context.Context, *CommandListRequest) (*CommandListResponse, error)
	CallCommand(context.Context, *CallCommandRequest) (*CallCommandResponse, error)
	DeleteCommand(context.Context, *DeleteCommandRequest) (*DeleteCommandResponse, error)
	CreateCommand(context.Context, *CreateCommandRequest) (*CreateCommandResponse, error)
	mustEmbedUnimplementedPiKvmAutomatorServer()
}

// UnimplementedPiKvmAutomatorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPiKvmAutomatorServer struct{}

func (UnimplementedPiKvmAutomatorServer) CommandList(context.Context, *CommandListRequest) (*CommandListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandList not implemented")
}
func (UnimplementedPiKvmAutomatorServer) CallCommand(context.Context, *CallCommandRequest) (*CallCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallCommand not implemented")
}
func (UnimplementedPiKvmAutomatorServer) DeleteCommand(context.Context, *DeleteCommandRequest) (*DeleteCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCommand not implemented")
}
func (UnimplementedPiKvmAutomatorServer) CreateCommand(context.Context, *CreateCommandRequest) (*CreateCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommand not implemented")
}
func (UnimplementedPiKvmAutomatorServer) mustEmbedUnimplementedPiKvmAutomatorServer() {}
func (UnimplementedPiKvmAutomatorServer) testEmbeddedByValue()                        {}

// UnsafePiKvmAutomatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PiKvmAutomatorServer will
// result in compilation errors.
type UnsafePiKvmAutomatorServer interface {
	mustEmbedUnimplementedPiKvmAutomatorServer()
}

func RegisterPiKvmAutomatorServer(s grpc.ServiceRegistrar, srv PiKvmAutomatorServer) {
	// If the following call pancis, it indicates UnimplementedPiKvmAutomatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PiKvmAutomator_ServiceDesc, srv)
}

func _PiKvmAutomator_CommandList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiKvmAutomatorServer).CommandList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PiKvmAutomator_CommandList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiKvmAutomatorServer).CommandList(ctx, req.(*CommandListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiKvmAutomator_CallCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiKvmAutomatorServer).CallCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PiKvmAutomator_CallCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiKvmAutomatorServer).CallCommand(ctx, req.(*CallCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiKvmAutomator_DeleteCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiKvmAutomatorServer).DeleteCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PiKvmAutomator_DeleteCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiKvmAutomatorServer).DeleteCommand(ctx, req.(*DeleteCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PiKvmAutomator_CreateCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PiKvmAutomatorServer).CreateCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PiKvmAutomator_CreateCommand_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PiKvmAutomatorServer).CreateCommand(ctx, req.(*CreateCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PiKvmAutomator_ServiceDesc is the grpc.ServiceDesc for PiKvmAutomator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PiKvmAutomator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pikvm_automator.PiKvmAutomator",
	HandlerType: (*PiKvmAutomatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommandList",
			Handler:    _PiKvmAutomator_CommandList_Handler,
		},
		{
			MethodName: "CallCommand",
			Handler:    _PiKvmAutomator_CallCommand_Handler,
		},
		{
			MethodName: "DeleteCommand",
			Handler:    _PiKvmAutomator_DeleteCommand_Handler,
		},
		{
			MethodName: "CreateCommand",
			Handler:    _PiKvmAutomator_CreateCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pikvm_automator.proto",
}
