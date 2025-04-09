package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"market/internal/database"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql      database.MysqlConfig
	CacheRedis cache.CacheConf
	Mongo      database.MongoConfig
}
