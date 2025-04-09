package config

import (
	"jobcenter/internal/database"
	"jobcenter/internal/logic"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Config struct {
	Okx        logic.OkxConfig
	Mongo      database.MongoConfig
	CacheRedis cache.CacheConf
}
