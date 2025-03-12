package handler

import (
	"ucenter-api/internal/midd"
	"ucenter-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	register := NewRegisterHandler(serverCtx)
	registerGroup := r.Group()
	registerGroup.Post("/uc/register/phone", register.Register)
	registerGroup.Post("/uc/mobile/code", register.SendCode)

	login := NewLoginHandler(serverCtx)
	loginGroup := r.Group()
	loginGroup.Post("/uc/login", login.Login)
	loginGroup.Post("/uc/check/login", login.CheckLogin)

	assetGroup := r.Group()
	assetGroup.Use(midd.Auth(serverCtx.Config.JWT.AccessSecret))
	asset := NewAssetHandler(serverCtx)
	assetGroup.Post("/uc/asset/wallet/:coinName", asset.FindWalletBySymbol)
	assetGroup.Post("/uc/asset/wallet", asset.FindWallet)
}
