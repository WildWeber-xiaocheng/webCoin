package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"
)

type MarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (l *MarketLogic) SymbolThumbTrend(req *types.MarketReq) (list []*types.CoinThumbResp, err error) {
	var thumbs []*market.CoinThumb
	thumb := l.svcCtx.Processor.GetThumb()
	isCache := false
	if thumb != nil {
		switch thumb.(type) {
		case []*market.CoinThumb:
			thumbs = thumb.([]*market.CoinThumb)
			isCache = true
		}
	}
	if !isCache {
		ctx, cancelFunc := context.WithTimeout(l.ctx, 10*time.Second)
		defer cancelFunc()
		symbolThumbRes, err := l.svcCtx.MarketRpc.FindSymbolThumbTrend(ctx,
			&market.MarketReq{
				Ip: req.Ip,
			})
		if err != nil {
			return nil, err
		}
		thumbs = symbolThumbRes.List
	}
	if err := copier.Copy(&list, thumbs); err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Trend == nil {
			v.Trend = []float64{}
		}
	}
	return
}

func (l *MarketLogic) SymbolThumb(req *types.MarketReq) (list []*types.CoinThumbResp, err error) {
	var thumbs []*market.CoinThumb
	thumb := l.svcCtx.Processor.GetThumb()
	if thumb != nil {
		switch thumb.(type) {
		case []*market.CoinThumb:
			thumbs = thumb.([]*market.CoinThumb)
		}
	}
	if err := copier.Copy(&list, thumbs); err != nil {
		return nil, err
	}
	return
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
