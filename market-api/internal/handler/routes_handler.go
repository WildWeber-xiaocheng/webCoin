package handler

import (
	"market-api/internal/svc"
)

func ExchangeRateHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	rate := NewExchangeRateHandler(serverCtx)
	rateGroup := r.Group()
	rateGroup.Post("/exchange-rate/usd/:unit", rate.GetUsdRate)
}
