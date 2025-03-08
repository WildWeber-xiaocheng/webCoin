package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"market-api/internal/svc"
	"market-api/internal/types"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *OrderLogic) History(req *types.OrderReq) (*types.HistoryKline, error) {
	//ctx, cancel := context.WithTimeout(l.ctx, 10*time.Second)
	//defer cancel()
	//return &types.HistoryKline{
	//	List: list,
	//}, nil
	return nil, nil
}

func (l *OrderLogic) Current(req *types.OrderReq) (*types.HistoryKline, error) {
	//ctx, cancel := context.WithTimeout(l.ctx, 10*time.Second)
	//defer cancel()
	//return &types.HistoryKline{
	//	List: list,
	//}, nil
	return nil, nil
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
