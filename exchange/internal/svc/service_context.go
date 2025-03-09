package svc

import (
	"exchange/internal/config"
	"exchange/internal/database"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/ucclient"
	"webCoin-common/msdb"
)

type ServiceContext struct {
	Config      config.Config
	Cache       cache.Cache
	Db          *msdb.MsDB
	MongoClient *database.MongoClient
	MemberRpc   ucclient.Member
	MarketRpc   mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(c.CacheRedis, nil, cache.NewStat("webCoin"), nil, func(o *cache.Options) {})
	return &ServiceContext{
		Config:      c,
		Cache:       redisCache,
		Db:          database.ConnMysql(c.Mysql.DataSource),
		MongoClient: database.ConnectMongo(c.Mongo),
		MemberRpc:   ucclient.NewMember(zrpc.MustNewClient(c.UCenterRpc)),
		MarketRpc:   mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
