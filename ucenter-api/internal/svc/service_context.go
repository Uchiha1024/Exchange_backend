package svc

import (
	"grpc-common/market/mclient"
	"grpc-common/ucenter/ucclient"
	"ucenter-api/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UCRegisterRpc ucclient.Register
	UCLoginRpc    ucclient.Login
	UCAssetRpc    ucclient.Asset
	MarketRpc     mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UCRegisterRpc: ucclient.NewRegister(zrpc.MustNewClient(c.UcenterRpc)),
		UCLoginRpc:    ucclient.NewLogin(zrpc.MustNewClient(c.UcenterRpc)),
		UCAssetRpc:    ucclient.NewAsset(zrpc.MustNewClient(c.UcenterRpc)),
		MarketRpc:     mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
