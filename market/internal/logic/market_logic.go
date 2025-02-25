package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market/internal/domain"
	"market/internal/svc"
	"math/rand"
	"time"
)

type MarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeCoinDomain *domain.ExchangeCoinDomain
}

func (l *MarketLogic) FindSymbolThumbTrend(req *market.MarketReq) (*market.SymbolThumbRes, error) {
	//查出所有的货币
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exchangeCoins := l.exchangeCoinDomain.FindVisible(ctx)
	coinThumbs := make([]*market.CoinThumb, len(exchangeCoins))
	for i, v := range exchangeCoins {
		trend := make([]float64, 0)
		for p := 0; p <= 24; p++ {
			trend = append(trend, rand.Float64())
		}
		ct := &market.CoinThumb{}
		ct.Symbol = v.Symbol
		ct.Trend = trend
		coinThumbs[i] = ct
	}
	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		exchangeCoinDomain: domain.NewExchangeCoinDomain(svcCtx.Db),
	}
}
