package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"ucenter/internal/database"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql       MysqlConfig
	CacheRedis  cache.CacheConf
	Captcha     CaptchaConf
	JWT         AuthConfig
	MarketRpc   zrpc.RpcClientConf
	ExchangeRpc zrpc.RpcClientConf
	Kafka       database.KafkaConfig
}

type MysqlConfig struct {
	DataSource string
}

type CaptchaConf struct {
	Vid string
	Key string
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}
