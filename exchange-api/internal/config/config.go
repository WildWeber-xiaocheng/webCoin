package config

import (
	"exchange-api/internal/database"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	MarketRpc zrpc.RpcClientConf
	Prefix    string
	KafKa     database.KafkaConfig
}
