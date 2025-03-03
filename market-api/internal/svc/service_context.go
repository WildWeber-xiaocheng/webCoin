package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"market-api/internal/config"
	"market-api/internal/database"
	"market-api/internal/processor"
	"market-api/internal/websocket"
)

type ServiceContext struct {
	Config          config.Config
	ExchangeRateRpc mclient.ExchangeRate
	MarketRpc       mclient.Market
	Processor       processor.Processor
}

func NewServiceContext(c config.Config, server *websocket.WebSocketServer) *ServiceContext {
	//初始化processor
	kafakaCli := database.NewKafkaClient(c.KafKa)
	market := mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc))
	defaultProcessor := processor.NewDefaultProcessor(kafakaCli)
	defaultProcessor.Init(market)
	defaultProcessor.AddHandler(processor.NewWebsocketHandler(server))

	return &ServiceContext{
		Config:          c,
		ExchangeRateRpc: mclient.NewExchangeRate(zrpc.MustNewClient(c.MarketRpc)),
		MarketRpc:       market,
		Processor:       defaultProcessor,
	}
}
