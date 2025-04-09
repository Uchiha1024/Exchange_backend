package svc

import (

	"ucenter-api/internal/config"
	"grpc-common/ucenter/ucclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UCRegisterRpc ucclient.Register
	UCLoginRpc    ucclient.Login
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UCRegisterRpc: ucclient.NewRegister(zrpc.MustNewClient(c.UcenterRpc)),
		UCLoginRpc:    ucclient.NewLogin(zrpc.MustNewClient(c.UcenterRpc)),
	}
}
