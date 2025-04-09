package svc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"jobcenter/internal/config"
	"jobcenter/internal/database"
)

type ServiceContext struct {
	Config      config.Config
	MongoClient *database.MongoClient
	Cache       cache.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化redis
	redisCache := cache.New(
		c.CacheRedis,
		nil,
		cache.NewStat("mscoin"),
		nil,
		func(o *cache.Options) {})

	return &ServiceContext{
		Config:      c,
		MongoClient: database.ConnectMongo(c.Mongo),
		Cache:       redisCache,
	}
}
