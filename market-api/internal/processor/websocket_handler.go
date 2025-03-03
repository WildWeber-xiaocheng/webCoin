package processor

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/model"
	"market-api/internal/websocket"
)

type WebsocketHandler struct {
	wsServer *websocket.WebSocketServer
}

func (w *WebsocketHandler) HandlerTrade(symbol string, data []byte) {
	//TODO implement me
	panic("implement me")
}

func (w *WebsocketHandler) HandlerKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) {
	logx.Info("======接收到数据,symbol=", symbol)
	thumb := thumbMap[symbol]
	if thumb == nil {
		thumb = kline.InitCoinThumb(symbol)
	}
	coinThumb := kline.ToCoinThumb(symbol, thumb)

	marshal, _ := json.Marshal(coinThumb) //转为json
	w.wsServer.BroadcastToNamespace("/", "/topic/market/thumb", string(marshal))
	logx.Info("=======接收到数据,kline=", string(marshal))
}

func NewWebsocketHandler(wsServer *websocket.WebSocketServer) *WebsocketHandler {
	return &WebsocketHandler{
		wsServer: wsServer,
	}
}
