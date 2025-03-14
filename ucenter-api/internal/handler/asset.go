package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"ucenter-api/internal/logic"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"
	common "webCoin-common"
	"webCoin-common/tools"
)

type AssetHandler struct {
	svcCtx *svc.ServiceContext
}

func (h *AssetHandler) FindWalletBySymbol(w http.ResponseWriter, r *http.Request) {
	// 1.获取参数
	var req types.AssetReq
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	// 2.获取用户钱包
	l := logic.NewAssetLogic(r.Context(), h.svcCtx)
	resp, err := l.FindWalletBySymbol(&req)
	// 3.处理返回数据
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *AssetHandler) FindWallet(w http.ResponseWriter, r *http.Request) {
	// 1.获取参数
	req := types.AssetReq{}
	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	// 2.获取用户钱包
	l := logic.NewAssetLogic(r.Context(), h.svcCtx)
	resp, err := l.FindWallet(&req)
	// 3.处理返回数据
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func (h *AssetHandler) ResetWalletAddress(w http.ResponseWriter, r *http.Request) {
	var req = types.AssetReq{}
	if err := httpx.ParseForm(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	ip := tools.GetRemoteClientIp(r)
	req.Ip = ip
	l := logic.NewAssetLogic(r.Context(), h.svcCtx)
	resp, err := l.ResetWalletAddress(&req)
	result := common.NewResult().Deal(resp, err)
	httpx.OkJsonCtx(r.Context(), w, result)
}

func NewAssetHandler(svcCtx *svc.ServiceContext) *AssetHandler {
	return &AssetHandler{svcCtx}
}
