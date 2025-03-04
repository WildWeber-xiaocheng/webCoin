// Code genemarketd by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: market.proto

package mclient

import (
	"context"
	"grpc-common/market/types/market"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	MarketReq      = market.MarketReq
	SymbolThumbRes = market.SymbolThumbRes
	ExchangeCoin   = market.ExchangeCoin

	Market interface {
		FindSymbolThumbTrend(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error)
		FindSymbolInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoin, error)
	}

	defaultMarket struct {
		cli zrpc.Client
	}
)

func NewMarket(cli zrpc.Client) Market {
	return &defaultMarket{
		cli: cli,
	}
}

func (m *defaultMarket) FindSymbolThumbTrend(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*SymbolThumbRes, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindSymbolThumbTrend(ctx, in, opts...)
}

func (m *defaultMarket) FindSymbolInfo(ctx context.Context, in *MarketReq, opts ...grpc.CallOption) (*ExchangeCoin, error) {
	client := market.NewMarketClient(m.cli.Conn())
	return client.FindSymbolInfo(ctx, in, opts...)
}
