package handler

import (
	"errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	common "webCoin-common"
	"webCoin-common/tools"
)

type LoginHandler struct {
	svcCtx *svc.ServiceContext
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginReq
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	result := common.NewResult()
	if req.Captcha == nil {
		httpx.OkJsonCtx(r.Context(), w, result.Deal(nil, errors.New("人机验证不通过")))
		return
	}

	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	l := logic.NewLoginLogic(r.Context(), h.svcCtx)
	resp, err := l.Login(&req)
	result = result.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *LoginHandler) CheckLogin(w http.ResponseWriter, r *http.Request) {
	result := common.NewResult()
	token := r.Header.Get("x-auth-token") //获取token
	l := logic.NewLoginLogic(r.Context(), h.svcCtx)
	resp, err := l.CheckLogin(token)
	result = result.Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func NewLoginHandler(svcCtx *svc.ServiceContext) *LoginHandler {
	return &LoginHandler{svcCtx}
}
