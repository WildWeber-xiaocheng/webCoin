package svc

import (
	"exchange-api/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
)

type ServiceContext struct {
	Config          config.Config
	ExchangeRateRpc mclient.ExchangeRate
	MarketRpc       mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	//初始化processor
	//kafaCli := database.NewKafkaClient(c.KafKa)
	market := mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc))
	return &ServiceContext{
		Config:          c,
		ExchangeRateRpc: mclient.NewExchangeRate(zrpc.MustNewClient(c.MarketRpc)),
		MarketRpc:       market,
	}
}
