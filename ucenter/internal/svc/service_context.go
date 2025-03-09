package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/exchange/eclient"
	"grpc-common/market/mclient"
	"ucenter/internal/config"
	"ucenter/internal/consumer"
	"ucenter/internal/database"
	"webCoin-common/msdb"
)

type ServiceContext struct {
	Config    config.Config
	Cache     cache.Cache
	Db        *msdb.MsDB
	MarketRpc mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(c.CacheRedis, nil, cache.NewStat("webCoin"),
		nil, func(o *cache.Options) {})
	mysql := database.ConnMysql(c.Mysql.DataSource)
	kafkaClient := database.NewKafkaClient(c.Kafka)
	kafkaClient.StartRead("add-exchange-order")
	order := eclient.NewOrder(zrpc.MustNewClient(c.ExchangeRpc))
	conf := c.CacheRedis[0].RedisConf
	newRedis := redis.MustNewRedis(conf)
	go consumer.ExchangeOrderAdd(newRedis, kafkaClient, order, mysql)
	return &ServiceContext{
		Config:    c,
		Cache:     redisCache,
		Db:        mysql,
		MarketRpc: mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
