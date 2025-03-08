package handler

import (
	"exchange-api/internal/svc"
)

func RegisterHandlers(r *Routers, serverCtx *svc.ServiceContext) {
	//如果要有中间件 怎么办？
	order := NewOrderHandler(serverCtx)
	orderGroup := r.Group()
	//历史委托订单：所有的订单
	orderGroup.Post("/order/history", order.History)
	//当前委托订单：正在交易的状态的订单
	orderGroup.Post("/order/current", order.Current)
}
