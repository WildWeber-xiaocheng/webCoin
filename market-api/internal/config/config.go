package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"market-api/internal/database"
)

type Config struct {
	rest.RestConf
	MarketRpc zrpc.RpcClientConf
	Prefix    string
	KafKa     database.KafkaConfig
}
