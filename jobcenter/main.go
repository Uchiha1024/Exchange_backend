package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"jobcenter/internal/config"
	"jobcenter/internal/svc"
	"jobcenter/internal/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	//日志的打印格式替换一下
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	t := task.NewTask(ctx)
	t.Run()
	//优雅退出
	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			log.Println("任务中心中断执行，开始clear资源")
			t.Stop()
			ctx.MongoClient.Disconnect()
		
			
		}
	}()
	t.StartBlocking()
}
