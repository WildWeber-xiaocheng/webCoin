package handler

import (
	"exchange-api/internal/svc"
	"net/http"
)

type OrderHandler struct {
	svcCtx *svc.ServiceContext
}

func (h OrderHandler) History(w http.ResponseWriter, r *http.Request) {

}

func (h OrderHandler) Current(w http.ResponseWriter, r *http.Request) {

}

func NewOrderHandler(svcCtx *svc.ServiceContext) *OrderHandler {
	return &OrderHandler{
		svcCtx: svcCtx,
	}
}
