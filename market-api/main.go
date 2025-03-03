package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/chain"
	"market-api/internal/config"
	"market-api/internal/handler"
	"market-api/internal/svc"
	"market-api/internal/websocket"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	//日志的格式修改
	logx.MustSetup(logx.LogConf{
		Stat:     false,
		Encoding: "plain",
	})
	var c config.Config
	conf.MustLoad(*configFile, &c)
	//使用websocket服务
	websocketServer := websocket.NewWebsocketServer("/socket.io")
	//解决跨域问题
	server := rest.MustNewServer(
		c.RestConf,
		rest.WithChain(chain.New(websocketServer.ServerHandler)), //将websocket服务注册到go-zero中
		rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
		}, nil, "http://localhost:8080"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c, websocketServer)
	//封装一层路由
	routers := handler.NewRouters(server, c.Prefix)
	handler.ExchangeRateHandlers(routers, ctx)

	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(websocketServer)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	group.Start()
}
