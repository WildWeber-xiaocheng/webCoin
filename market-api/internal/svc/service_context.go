package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"market-api/internal/config"
	"market-api/internal/database"
	"market-api/internal/processor"
)

type ServiceContext struct {
	Config          config.Config
	ExchangeRateRpc mclient.ExchangeRate
	MarketRpc       mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	//初始化processor
	kafakaCli := database.NewKafkaClient(c.KafKa)
	defaultProcessor := processor.NewDefaultProcessor(kafakaCli)
	defaultProcessor.Init()
	defaultProcessor.AddHandler(processor.NewWebsocketHandler())
	return &ServiceContext{
		Config:          c,
		ExchangeRateRpc: mclient.NewExchangeRate(zrpc.MustNewClient(c.MarketRpc)),
		MarketRpc:       mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
