package svc

import (
	"exchange-api/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
	"grpc-common/exchange/eclient"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc eclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	order := eclient.NewOrder(zrpc.MustNewClient(c.ExchangeRpc))
	return &ServiceContext{
		Config:   c,
		OrderRpc: order,
	}
}
