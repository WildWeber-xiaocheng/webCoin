package handler

import (
	"market-api/internal/svc"
)

func ExchangeRateHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	rate := NewExchangeRateHandler(serverCtx)
	rateGroup := r.Group()
	rateGroup.Post("/exchange-rate/usd/:unit", rate.GetUsdRate)
	market := NewMarketHandler(serverCtx)
	marketGroup := r.Group()
	marketGroup.Post("/symbol-thumb-trend", market.SymbolThumbTrend)

	//要注册websocket路径
	wsGroup := r.Group()
	wsGroup.GetNoPrefix("/socket.io", nil)
	wsGroup.PostNoprefix("/socket.io", nil)
}
