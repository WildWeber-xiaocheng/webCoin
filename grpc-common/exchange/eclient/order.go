// Code geneorderd by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: order.proto

package eclient

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"grpc-common/exchange/types/order"
)

type (
	OrderReq = order.OrderReq
	OrderRes = order.OrderRes

	Order interface {
		FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
		FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func (d *defaultOrder) FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.FindOrderHistory(ctx, in, opts...)
}

func (d *defaultOrder) FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.FindOrderCurrent(ctx, in, opts...)
}

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}
