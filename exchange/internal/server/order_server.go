// Code genemarketd by goctl. DO NOT EDIT.
// goctl 1.7.6
// Source: register.proto

package server

import (
	"exchange/internal/svc"
	"grpc-common/market/types/market"
)

type MarketServer struct {
	svcCtx *svc.ServiceContext
	market.UnimplementedMarketServer
}

//func (e *MarketServer) FindSymbolThumbTrend(ctx context.Context, req *market.MarketReq) (*market.SymbolThumbRes, error) {
//	l := logic.NewMarketLogic(ctx, e.svcCtx)
//	return l.FindSymbolThumbTrend(req)
//}

func NewMarketServer(svcCtx *svc.ServiceContext) *MarketServer {
	return &MarketServer{
		svcCtx: svcCtx,
	}
}
