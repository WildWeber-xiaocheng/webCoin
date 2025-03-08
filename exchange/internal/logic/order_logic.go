package logic

import (
	"context"
	"exchange/internal/domain"
	"exchange/internal/svc"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeOrderDomain *domain.ExchangeOrderDomain
}

func (l *ExchangeOrderLogic) FindOrderHistory(req *order.OrderReq) (*order.OrderRes, error) {
	exchangeOrders, total, err := l.exchangeOrderDomain.FindOrderHistory(
		l.ctx,
		req.Symbol,
		req.Page,
		req.PageSize,
		req.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) FindOrderCurrent(req *order.OrderReq) (*order.OrderRes, error) {
	exchangeOrders, total, err := l.exchangeOrderDomain.FindOrderCurrent(
		l.ctx,
		req.Symbol,
		req.Page,
		req.PageSize,
		req.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	return &ExchangeOrderLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		exchangeOrderDomain: domain.NewExchangeOrderDomain(svcCtx.Db),
	}
}
