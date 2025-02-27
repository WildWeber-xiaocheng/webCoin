package domain

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market/internal/dao"
	"market/internal/database"
	"market/internal/model"
	"market/internal/repo"
	"time"
	"webCoin-common/tools"
)

type MarketDomain struct {
	klineRepo repo.KlineRepo
}

func (d *MarketDomain) SymbolThumbTrend(coins []*model.ExchangeCoin) []*market.CoinThumb {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	list := make([]*market.CoinThumb, len(coins))
	to := time.Now().UnixMilli()
	from := tools.ZeroTime()
	for i, v := range coins {
		klines, err := d.klineRepo.FindBySymbolTime(ctx, v.Symbol, "1H", from, to)
		if err != nil {
			logx.Error(err)
			list[i] = model.DefaultCoinThumb(v.Symbol)
			continue
		}
		klinesLength := len(klines)
		if klinesLength <= 0 {
			list[i] = model.DefaultCoinThumb(v.Symbol)
			continue
		}
		//因为FindBySymbolTime是降序排列，即klines[0]是最新数据,klines[length-1]是当天0点的数据
		trend := make([]float64, len(klines))
		var high float64 = 0
		low := klines[0].LowestPrice
		var volumes float64 = 0
		var turnover float64 = 0
		for i, v := range klines {
			trend[i] = v.ClosePrice
			if v.HighestPrice > high {
				high = v.HighestPrice
			}
			if v.LowestPrice < low {
				low = v.LowestPrice
			}
			volumes += v.Volume
			turnover += v.Turnover
		}
		newKline := klines[0]
		oldKline := klines[len(klines)-1]
		ct := newKline.ToCoinThumb(v.Symbol, oldKline)
		ct.High = high
		ct.Low = low
		ct.Volume = volumes
		ct.Turnover = turnover
		ct.Trend = trend
		list[i] = ct
	}
	return list
}

func NewMarketDomain(mongoClient *database.MongoClient) *MarketDomain {
	return &MarketDomain{
		klineRepo: dao.NewKlineDao(mongoClient.Db),
	}
}
