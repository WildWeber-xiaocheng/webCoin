package logic

import (
	"context"
	"exchange/internal/domain"
	"exchange/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeOrderDomain *domain.ExchangeOrderDomain
}

func (l *ExchangeOrderLogic) HistoryKline(req *market.MarketReq) (*market.HistoryRes, error) {
	return &market.HistoryRes{}, nil
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	return &ExchangeOrderLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		exchangeOrderDomain: domain.NewExchangeOrderDomain(svcCtx.Db),
	}
}
