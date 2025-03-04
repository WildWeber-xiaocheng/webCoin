package handler

import (
	"market-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	//如果要有中间件 怎么办？
	rate := NewExchangeRateHandler(serverCtx)
	rateGroup := r.Group()
	rateGroup.Post("/exchange-rate/usd/:unit", rate.GetUsdRate)
	market := NewMarketHandler(serverCtx)
	marketGroup := r.Group()
	marketGroup.Post("/symbol-thumb-trend", market.SymbolThumbTrend)
	marketGroup.Post("/symbol-thumb", market.SymbolThumb)
	marketGroup.Post("/symbol-info", market.SymbolInfo)

	wsGroup := r.Group()
	wsGroup.GetNoPrefix("/socket.io", nil)
	wsGroup.PostNoPrefix("/socket.io", nil)
}
