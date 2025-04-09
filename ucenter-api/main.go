package main

import (
	"flag"
	"fmt"
	"net/http"

	"ucenter-api/internal/config"
	"ucenter-api/internal/handler"
	"ucenter-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	// 日志打印格式调整
	logx.MustSetup(logx.LogConf{
		Stat:     false,
		Encoding: "plain",
	})

	var c config.Config
	conf.MustLoad(*configFile, &c)
	// 解决跨域问题
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
	}, nil, "http://localhost:8080"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	router := handler.NewRoutes(server)
	handler.RegisterHandlers(router, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
