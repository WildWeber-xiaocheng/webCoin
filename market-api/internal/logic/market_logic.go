package logic

import (
	"context"
	"errors"
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
	var coinThumbs []*market.CoinThumb
	thumbs := l.svcCtx.Processor.GetThumb()
	isCache := false
	if thumbs != nil {
		switch thumbs.(type) {
		case []*market.CoinThumb:
			coinThumbs = thumbs.([]*market.CoinThumb)
			isCache = true
		}
	}
	if !isCache {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		thumbResp, err := l.svcCtx.MarketRpc.FindSymbolThumbTrend(ctx, &market.MarketReq{
			Ip: req.Ip,
		})
		if err != nil {
			return nil, err
		}
		coinThumbs = thumbResp.List
	}
	if err := copier.Copy(&list, coinThumbs); err != nil {
		return nil, errors.New("数据格式有误")
	}
	for _, v := range list {
		if v.Trend == nil {
			v.Trend = []float64{}
		}
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
