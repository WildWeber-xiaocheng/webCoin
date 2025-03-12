package processor

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/types/market"
	"market-api/internal/model"
	"market-api/internal/ws"
)

type WebsocketHandler struct {
	wsServer *ws.WebsocketServer
}

func (w *WebsocketHandler) HandleTradePlate(symbol string, tp *model.TradePlateResult) {
	marshal, _ := json.Marshal(tp)
	logx.Info("====买卖盘通知:", symbol, tp.Direction, ":", fmt.Sprintf("%d", len(tp.Items)))
	w.wsServer.BroadcastToNamespace("/", "/topic/market/trade-plate/"+symbol, string(marshal))
}

func (w *WebsocketHandler) HandleTrade(symbol string, data []byte) {
	//本方法没有实现
	//本方法原本是在订单交易完成后进入这个这个函数，将订单的数据作为k线的一部分
	//由于k线的数据来源于第三方，为了保证k线数据的真实性，这里的虚拟交易就不作为k线数据了
	//TODO implement me
	panic("implement me")
}

func (w *WebsocketHandler) HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) {
	logx.Info("================WebsocketHandler Start=======================")
	logx.Info("symbol:", symbol)
	thumb := thumbMap[symbol]
	if thumb == nil {
		thumb = kline.InitCoinThumb(symbol)
	}
	coinThumb := kline.ToCoinThumb(symbol, thumb)
	result := &model.CoinThumb{}
	copier.Copy(result, coinThumb)
	marshal, _ := json.Marshal(result)
	w.wsServer.BroadcastToNamespace("/", "/topic/market/thumb", string(marshal))

	bytes, _ := json.Marshal(kline)
	w.wsServer.BroadcastToNamespace("/", "/topic/market/kline/"+symbol, string(bytes))

	logx.Info("marshal:", marshal)
	logx.Info("================WebsocketHandler End=======================")
}

func NewWebsocketHandler(wsServer *ws.WebsocketServer) *WebsocketHandler {
	return &WebsocketHandler{
		wsServer: wsServer,
	}
}
