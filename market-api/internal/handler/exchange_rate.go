package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"market-api/internal/logic"
	"market-api/internal/svc"
	"market-api/internal/types"
	"net/http"
	common "webCoin-common"
	"webCoin-common/tools"
)

type ExchangeRateHandler struct {
	svcCtx *svc.ServiceContext
}

func NewExchangeRateHandler(svcCtx *svc.ServiceContext) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		svcCtx: svcCtx,
	}
}

func (h *ExchangeRateHandler) GetUsdRate(w http.ResponseWriter, r *http.Request) {
	var req types.RateRequest
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	//获取ip
	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	l := logic.NewExchangeRateLogic(r.Context(), h.svcCtx)
	resp, err := l.GetUsdCnyRate(&req)
	result := common.NewResult().Deal(resp.Rate, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}
