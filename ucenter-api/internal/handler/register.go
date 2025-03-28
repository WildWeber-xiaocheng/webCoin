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

type RegisterHandler struct {
	svcCtx *svc.ServiceContext
}

func NewRegisterHandler(svcCtx *svc.ServiceContext) *RegisterHandler {
	return &RegisterHandler{
		svcCtx: svcCtx,
	}
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.Request
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	result := common.NewResult()
	if req.Captcha == nil {
		httpx.OkJsonCtx(r.Context(), w, result.Deal(nil, errors.New("人机验证不通过")))
		return
	}

	//获取ip
	req.Ip = tools.GetRemoteClientIp(r)

	l := logic.NewRegisterLogic(r.Context(), h.svcCtx)
	resp, err := l.Register(&req)
	newResult := result.Deal(resp, err)
	//成功与否状态码都返回200
	httpx.OkJsonCtx(r.Context(), w, newResult)
}

func (h *RegisterHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	var req types.CodeRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	l := logic.NewRegisterLogic(r.Context(), h.svcCtx)
	resp, err := l.SendCode(&req)
	result := common.NewResult().Deal(resp, err)
	//成功与否状态码都返回200
	httpx.OkJsonCtx(r.Context(), w, result)
}
