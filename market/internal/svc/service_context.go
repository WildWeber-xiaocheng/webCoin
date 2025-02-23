package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"market/internal/config"
	"market/internal/database"
	"webCoin-common/msdb"
)

type ServiceContext struct {
	Config config.Config
	Cache  cache.Cache
	Db     *msdb.MsDB
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(c.CacheRedis, nil, cache.NewStat("webCoin"), nil, func(o *cache.Options) {})
	return &ServiceContext{
		Config: c,
		Cache:  redisCache,
		Db:     database.ConnMysql(c.Mysql.DataSource),
	}
}
