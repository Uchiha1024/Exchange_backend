package svc

import (
	"mscoin-common/msdb"
	"ucenter/internal/config"
	"ucenter/internal/database"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

type ServiceContext struct {
	Config config.Config
	Redis  cache.Cache
	Db     *msdb.MsDB
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(
		c.CacheRedis,
		nil,
		cache.NewStat("mscoin"),
		nil,
		func(o *cache.Options) {})
	return &ServiceContext{
		Config: c,
		Redis:redisCache,
		Db: database.ConnMysql(c.Mysql.DataSource),

	}
}
