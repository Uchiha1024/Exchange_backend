package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UcenterRpc zrpc.RpcClientConf
	JWT        AuthConfig
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}
