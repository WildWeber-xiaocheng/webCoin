package config

import (
	"exchange-api/internal/database"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ExchangeRpc zrpc.RpcClientConf
	Prefix      string
	KafKa       database.KafkaConfig
	JWT         AuthConfig
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}
