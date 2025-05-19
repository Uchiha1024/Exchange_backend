package main

import (
	"flag"
	"fmt"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/login"
	"grpc-common/ucenter/types/member"
	"grpc-common/ucenter/types/register"
	"grpc-common/ucenter/types/withdraw"
	"ucenter/internal/config"
	"ucenter/internal/server"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	// 添加这行来禁用缓存统计功能
	// stat.SetReporter(nil)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		register.RegisterRegisterServer(grpcServer, server.NewRegisterServer(ctx))
		login.RegisterLoginServer(grpcServer, server.NewLoginServer(ctx))
		asset.RegisterAssetServer(grpcServer, server.NewAssetServer(ctx))
		member.RegisterMemberServer(grpcServer, server.NewMemberServer(ctx))
		withdraw.RegisterWithdrawServer(grpcServer, server.NewWithdrawServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
