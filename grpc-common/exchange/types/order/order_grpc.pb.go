// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: order.proto

package order

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

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
	FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	out := new(OrderRes)
	err := c.cc.Invoke(ctx, "/order.Order/FindOrderHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	out := new(OrderRes)
	err := c.cc.Invoke(ctx, "/order.Order/FindOrderCurrent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility
type OrderServer interface {
	FindOrderHistory(context.Context, *OrderReq) (*OrderRes, error)
	FindOrderCurrent(context.Context, *OrderReq) (*OrderRes, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServer struct {
}

func (UnimplementedOrderServer) FindOrderHistory(context.Context, *OrderReq) (*OrderRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOrderHistory not implemented")
}
func (UnimplementedOrderServer) FindOrderCurrent(context.Context, *OrderReq) (*OrderRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindOrderCurrent not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_FindOrderHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).FindOrderHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.Order/FindOrderHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).FindOrderHistory(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_FindOrderCurrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).FindOrderCurrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/order.Order/FindOrderCurrent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).FindOrderCurrent(ctx, req.(*OrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindOrderHistory",
			Handler:    _Order_FindOrderHistory_Handler,
		},
		{
			MethodName: "FindOrderCurrent",
			Handler:    _Order_FindOrderCurrent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
