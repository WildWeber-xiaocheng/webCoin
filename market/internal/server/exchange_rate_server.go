// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: register.proto

package server

import (
	"market/internal/svc"
)

type ExchangeRateServer struct {
	svcCtx *svc.ServiceContext
	//register.UnimplementedExchangeRateServer
}

func NewExchangeRateServer(svcCtx *svc.ServiceContext) *ExchangeRateServer {
	return &ExchangeRateServer{
		svcCtx: svcCtx,
	}
}
