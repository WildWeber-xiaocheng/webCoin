package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	common "webCoin-common"
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
	//if err := httpx.Parse(r, &req); err != nil {
	//	httpx.ErrorCtx(r.Context(), w, err)
	//	return
	//}

	l := logic.NewRegisterLogic(r.Context(), h.svcCtx)
	resp, err := l.Register(&req)
	result := common.NewResult().Deal(resp, err)
	//成功与否状态码都返回200
	httpx.OkJsonCtx(r.Context(), w, result)
}
