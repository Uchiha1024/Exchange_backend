package config

import (
	"jobcenter/internal/database"
	"jobcenter/internal/logic"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Okx        logic.OkxConfig
	Mongo      database.MongoConfig
	CacheRedis cache.CacheConf
	Kafka      database.KafkaConfig
	UCenterRpc zrpc.RpcClientConf
	Bitcoin    logic.BitCoinConfig
}
