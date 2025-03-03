package processor

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"market-api/internal/model"
)

type WebsocketHandler struct {
}

func (w *WebsocketHandler) HandlerTrade(symbol string, data []byte) {
	//TODO implement me
	panic("implement me")
}

func (w *WebsocketHandler) HandlerKLine(symbol string, kline *model.Kline) {
	logx.Info("======接收到数据,symbol=", symbol)
	marshal, _ := json.Marshal(kline)
	logx.Info("=======接收到数据,kline=", string(marshal))
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}
