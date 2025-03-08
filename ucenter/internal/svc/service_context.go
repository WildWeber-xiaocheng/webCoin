package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"ucenter/internal/config"
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
	redisCache := cache.New(c.CacheRedis, nil, cache.NewStat("webCoin"), nil, func(o *cache.Options) {})
	return &ServiceContext{
		Config:    c,
		Cache:     redisCache,
		Db:        database.ConnMysql(c.Mysql.DataSource),
		MarketRpc: mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
