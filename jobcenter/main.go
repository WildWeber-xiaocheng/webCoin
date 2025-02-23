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
	//读配置
	var c config.Config
	conf.MustLoad(*configFile, &c)
	//
	ctx := svc.NewServiceContext(c)
	//开始任务
	t := task.NewTask(ctx)
	t.Run()
	//优雅退出
	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			log.Println("监听到中断信号，终止程序。任务中心中断，开始clear资源")
			t.Stop()
		}
	}()
	t.StartBlocking()
}
