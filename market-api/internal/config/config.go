package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"market-api/internal/database"
)

type Config struct {
	rest.RestConf
	Prefix    string
	MarketRpc zrpc.RpcClientConf
	Kafka     database.KafkaConfig
}
