package svc

import (
	"grpc-common/market/mclient"
	"mscoin-common/msdb"
	"ucenter/internal/config"
	"ucenter/internal/database"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Redis     cache.Cache
	Db        *msdb.MsDB
	MarketRpc mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.New(
		c.CacheRedis,
		nil,
		cache.NewStat("mscoin"),
		nil,
		func(o *cache.Options) {})
	return &ServiceContext{
		Config:    c,
		Redis:     redisCache,
		Db:        database.ConnMysql(c.Mysql.DataSource),
		MarketRpc: mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
