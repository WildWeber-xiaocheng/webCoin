package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market/internal/domain"
	"market/internal/svc"
	"time"
)

type MarketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeCoinDomain *domain.ExchangeCoinDomain
	marketDomain       *domain.MarketDomain
}

func (l *MarketLogic) FindSymbolThumbTrend(req *market.MarketReq) (*market.SymbolThumbRes, error) {
	//查出所有的货币
	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()
	exchangeCoins := l.exchangeCoinDomain.FindVisible(ctx)
	// 查询mongo中相应的数据
	// 这里查询1h间隔的数据，可以根据时间来进行查询，前端显示的24h范围的数据，这里简化为当天的价格变化趋势
	coinThumbs := l.marketDomain.SymbolThumbTrend(exchangeCoins)
	//测试数据，测试接口是否正常
	//coinThumbs := make([]*market.CoinThumb, len(exchangeCoins))
	//for i, v := range exchangeCoins {
	//	trend := make([]float64, 0)
	//	for p := 0; p <= 24; p++ {
	//		trend = append(trend, rand.Float64())
	//	}
	//	ct := &market.CoinThumb{}
	//	ct.Symbol = v.Symbol
	//	ct.Trend = trend
	//	coinThumbs[i] = ct
	//}
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
		marketDomain:       domain.NewMarketDomain(svcCtx.MongoClient),
	}
}
