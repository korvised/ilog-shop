// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: modules/item/itemPb/itemPb.proto

package ilog_shop

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

// ItemGrpcServiceClient is the client API for ItemGrpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ItemGrpcServiceClient interface {
	FindItemInIds(ctx context.Context, in *FindItemsInIdsReq, opts ...grpc.CallOption) (*FindItemsInIdsRes, error)
}

type itemGrpcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewItemGrpcServiceClient(cc grpc.ClientConnInterface) ItemGrpcServiceClient {
	return &itemGrpcServiceClient{cc}
}

func (c *itemGrpcServiceClient) FindItemInIds(ctx context.Context, in *FindItemsInIdsReq, opts ...grpc.CallOption) (*FindItemsInIdsRes, error) {
	out := new(FindItemsInIdsRes)
	err := c.cc.Invoke(ctx, "/ItemGrpcService/FindItemInIds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemGrpcServiceServer is the server API for ItemGrpcService service.
// All implementations must embed UnimplementedItemGrpcServiceServer
// for forward compatibility
type ItemGrpcServiceServer interface {
	FindItemInIds(context.Context, *FindItemsInIdsReq) (*FindItemsInIdsRes, error)
	mustEmbedUnimplementedItemGrpcServiceServer()
}

// UnimplementedItemGrpcServiceServer must be embedded to have forward compatible implementations.
type UnimplementedItemGrpcServiceServer struct {
}

func (UnimplementedItemGrpcServiceServer) FindItemInIds(context.Context, *FindItemsInIdsReq) (*FindItemsInIdsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindItemInIds not implemented")
}
func (UnimplementedItemGrpcServiceServer) mustEmbedUnimplementedItemGrpcServiceServer() {}

// UnsafeItemGrpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemGrpcServiceServer will
// result in compilation errors.
type UnsafeItemGrpcServiceServer interface {
	mustEmbedUnimplementedItemGrpcServiceServer()
}

func RegisterItemGrpcServiceServer(s grpc.ServiceRegistrar, srv ItemGrpcServiceServer) {
	s.RegisterService(&ItemGrpcService_ServiceDesc, srv)
}

func _ItemGrpcService_FindItemInIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindItemsInIdsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemGrpcServiceServer).FindItemInIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ItemGrpcService/FindItemInIds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemGrpcServiceServer).FindItemInIds(ctx, req.(*FindItemsInIdsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemGrpcService_ServiceDesc is the grpc.ServiceDesc for ItemGrpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemGrpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ItemGrpcService",
	HandlerType: (*ItemGrpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindItemInIds",
			Handler:    _ItemGrpcService_FindItemInIds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "modules/item/itemPb/itemPb.proto",
}
