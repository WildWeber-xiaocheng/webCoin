package logic

import (
	"context"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"webCoin-common/pages"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *OrderLogic) History(req *types.ExchangeReq) (*pages.PageResult, error) {
	value := l.ctx.Value("userId").(int64)
	history, err := l.svcCtx.OrderRpc.FindOrderHistory(l.ctx, &order.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   value,
	})
	if err != nil {
		return nil, err
	}
	list := history.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	//构建分页查询结果
	return pages.New(b, req.PageNo, req.PageSize, history.Total), nil
}

func (l *OrderLogic) Current(req *types.ExchangeReq) (*pages.PageResult, error) {
	value := l.ctx.Value("userId").(int64)
	history, err := l.svcCtx.OrderRpc.FindOrderCurrent(l.ctx, &order.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   value,
	})
	if err != nil {
		return nil, err
	}
	list := history.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, req.PageNo, req.PageSize, history.Total), nil
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
