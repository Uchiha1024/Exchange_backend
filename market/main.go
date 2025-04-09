package main

import (
	"flag"
	"fmt"
	"grpc-common/market/types/rate"
	"market/internal/config"
	"market/internal/server"
	"market/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	//日志的打印格式替换一下
	logx.MustSetup(logx.LogConf{Stat: false, Encoding: "plain"})
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rate.RegisterExchangeRateServer(grpcServer, server.NewExchangeRateServer(ctx))
		// market.RegisterMarketServer(grpcServer, server.NewMarketServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
