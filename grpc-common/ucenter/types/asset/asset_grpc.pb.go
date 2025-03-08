// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: asset.proto

package asset

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

// AssetClient is the client API for Asset service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssetClient interface {
	FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error)
}

type assetClient struct {
	cc grpc.ClientConnInterface
}

func NewAssetClient(cc grpc.ClientConnInterface) AssetClient {
	return &assetClient{cc}
}

func (c *assetClient) FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error) {
	out := new(MemberWallet)
	err := c.cc.Invoke(ctx, "/asset.Asset/findWalletBySymbol", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssetServer is the server API for Asset service.
// All implementations must embed UnimplementedAssetServer
// for forward compatibility
type AssetServer interface {
	FindWalletBySymbol(context.Context, *AssetReq) (*MemberWallet, error)
	mustEmbedUnimplementedAssetServer()
}

// UnimplementedAssetServer must be embedded to have forward compatible implementations.
type UnimplementedAssetServer struct {
}

func (UnimplementedAssetServer) FindWalletBySymbol(context.Context, *AssetReq) (*MemberWallet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindWalletBySymbol not implemented")
}
func (UnimplementedAssetServer) mustEmbedUnimplementedAssetServer() {}

// UnsafeAssetServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssetServer will
// result in compilation errors.
type UnsafeAssetServer interface {
	mustEmbedUnimplementedAssetServer()
}

func RegisterAssetServer(s grpc.ServiceRegistrar, srv AssetServer) {
	s.RegisterService(&Asset_ServiceDesc, srv)
}

func _Asset_FindWalletBySymbol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetServer).FindWalletBySymbol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/asset.Asset/findWalletBySymbol",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetServer).FindWalletBySymbol(ctx, req.(*AssetReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Asset_ServiceDesc is the grpc.ServiceDesc for Asset service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Asset_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "asset.Asset",
	HandlerType: (*AssetServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "findWalletBySymbol",
			Handler:    _Asset_FindWalletBySymbol_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "asset.proto",
}
