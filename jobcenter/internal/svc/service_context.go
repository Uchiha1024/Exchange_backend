package svc

import (
	"grpc-common/ucenter/ucclient"
	"jobcenter/internal/config"
	"jobcenter/internal/database"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	MongoClient    *database.MongoClient
	Cache          cache.Cache
	KafkaClient    *database.KafkaClient
	AssetRpc       ucclient.Asset
	BitCoinAddress string
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化redis
	redisCache := cache.New(
		c.CacheRedis,
		nil,
		cache.NewStat("mscoin"),
		nil,
		func(o *cache.Options) {})
	// 初始化 kafka
	client := database.NewKafkaClient(c.Kafka)
	client.StartWrite()

	return &ServiceContext{
		Config:         c,
		MongoClient:    database.ConnectMongo(c.Mongo),
		Cache:          redisCache,
		KafkaClient:    client,
		AssetRpc:       ucclient.NewAsset(zrpc.MustNewClient(c.UCenterRpc)),
		BitCoinAddress: c.Bitcoin.Address,
	}
}
