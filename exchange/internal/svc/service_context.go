package svc

import (
	"exchange/internal/config"
	"exchange/internal/consumer"
	"exchange/internal/database"
	"exchange/internal/processor"
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
	AssetRpc    ucclient.Asset
	KafkaClient *database.KafkaClient
}

func (sc *ServiceContext) init() {
	factory := processor.NewCoinTradeFactory()
	factory.Init(sc.MarketRpc, sc.KafkaClient, sc.Db)
	kafkaConsumer := consumer.NewKafkaConsumer(sc.KafkaClient, factory, sc.Db)
	kafkaConsumer.Run()
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(c.CacheRedis, nil, cache.NewStat("webCoin"), nil, func(o *cache.Options) {})
	kafkaClient := database.NewKafkaClient(c.Kafka)
	s := &ServiceContext{
		Config:      c,
		Cache:       redisCache,
		Db:          database.ConnMysql(c.Mysql.DataSource),
		MongoClient: database.ConnectMongo(c.Mongo),
		MemberRpc:   ucclient.NewMember(zrpc.MustNewClient(c.UCenterRpc)),
		MarketRpc:   mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
		AssetRpc:    ucclient.NewAsset(zrpc.MustNewClient(c.UCenterRpc)),
		KafkaClient: kafkaClient,
	}
	s.init()
	return s
}
