package websocket

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

const ROOM = "market"

type WebSocketServer struct {
	server *socketio.Server
	path   string
}

func (ws *WebSocketServer) Start() {
	logx.Info("============socketIO启动================")
	ws.server.Serve()
}

func (ws *WebSocketServer) Stop() {
	logx.Info("============socketIO关闭================")
	ws.server.Close()
}

// 跨域
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func (ws *WebSocketServer) BroadcastToNamespace(path string, event string, data any) {
	go func() {
		ws.server.BroadcastToRoom(path, ROOM, event, data)
	}()
}

func NewWebsocketServer(path string) *WebSocketServer {
	//NewServer传一个option，解决跨域问题
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		conn.Join(ROOM)
		logx.Info("connected:", conn.ID())
		return nil
	})
	return &WebSocketServer{
		server: server,
		path:   path,
	}
}

func (ws *WebSocketServer) ServerHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logx.Info("=============", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, ws.path) {
			//进行处理
			ws.server.ServeHTTP(w, r)
		} else { //放行
			next.ServeHTTP(w, r)
		}
	})
}
