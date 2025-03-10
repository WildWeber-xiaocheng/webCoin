// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: register.proto

package ucclient

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"grpc-common/ucenter/types/asset"
)

type (
	AssetReq     = asset.AssetReq
	MemberWallet = asset.MemberWallet

	Asset interface {
		FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error)
	}

	defaultAsset struct {
		cli zrpc.Client
	}
)

func NewAsset(cli zrpc.Client) Asset {
	return &defaultAsset{
		cli: cli,
	}
}

func (m *defaultAsset) FindWalletBySymbol(ctx context.Context, in *AssetReq, opts ...grpc.CallOption) (*MemberWallet, error) {
	client := asset.NewAssetClient(m.cli.Conn())
	return client.FindWalletBySymbol(ctx, in, opts...)
}
