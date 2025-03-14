package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/ucenter/ucclient"
	"jobcenter/internal/config"
	"jobcenter/internal/database"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *database.MongoClient
	KafkaClient *database.KafkaClient
	Cache       cache.Cache
	AssetRpc    ucclient.Asset
}

func NewServiceContext(c config.Config) *ServiceContext {
	kafkaClient := database.NewKafkaClient(c.Kafka)
	kafkaClient.StartWrite()
	redisCache := cache.New(
		c.CacheRedis,
		nil,
		cache.NewStat("webCoin"),
		nil,
		func(o *cache.Options) {})
	return &ServiceContext{
		Config:      c,
		MongoClient: database.ConnectMongo(c.Mongo),
		KafkaClient: kafkaClient,
		Cache:       redisCache,
		AssetRpc:    ucclient.NewAsset(zrpc.MustNewClient(c.UCenterRpc)),
	}
}
