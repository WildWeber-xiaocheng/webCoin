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
}
