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
	"market-api/internal/ws"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	//日志的打印格式替换一下
	logx.MustSetup(logx.LogConf{Stat: false, Encoding: "plain"})
	var c config.Config
	conf.MustLoad(*configFile, &c)
	wsServer := ws.NewWebsocketServer("/socket.io")
	fmt.Printf("wsServer: %+v\n", wsServer)
	fmt.Printf("Type of ServerHandler: %T\n", wsServer.ServerHandler)
	server := rest.MustNewServer(
		c.RestConf,
		//rest.WithChain(chain.New(wsServer.ServerHandler)),
		rest.WithChain(chain.New(func(next http.Handler) http.Handler {
			fmt.Println("ServerHandler middleware is being executed")
			return wsServer.ServerHandler(next)
		})),
		//rest.WithRouter(自定义的路由实现) zero框架就会走你的路由
		rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
		}, nil, "http://localhost:8080"))
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodGet,
	//	Path:    "/socket.io/",
	//	Handler: wsServer.ServerHandler.ServeHTTP,
	//})
	//server.AddRoute(rest.Route{
	//	Method: http.MethodGet,
	//	Path:   "/socket.io/",
	//	Handler: func(w http.ResponseWriter, r *http.Request) {
	//		wsServer.ServeHTTP(w, r)
	//	},
	//})
	defer server.Stop()

	ctx := svc.NewServiceContext(c, wsServer)
	router := handler.NewRouters(server, c.Prefix)
	handler.RegisterHandlers(router, ctx)

	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(wsServer)
	//fmt.Printf("Service context: %+v\n", ctx)
	//fmt.Printf("Router: %+v\n", router)
	//fmt.Printf("group: %+v\n", group)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	group.Start()
}
