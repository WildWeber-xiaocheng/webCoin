package handler

import (
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
}
