package repo

import (
	"context"
	"exchange/internal/model"
	"webCoin-common/msdb"
)

type ExchangeOrderRepo interface {
	FindOrderHistory(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error)
	FindOrderCurrent(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error)
	FindCurrentTradingCount(ctx context.Context, userId int64, symbol string, direction int) (int64, error)
	Save(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder) error
	FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error)
	UpdateOrderStatusCancel(ctx context.Context, orderId string) error
	UpdateOrderStatusTrading(ctx context.Context, orderId string) error
	FindOrderListBySymbol(ctx context.Context, symbol string, status int) ([]*model.ExchangeOrder, error)
	UpdateOrderComplete(ctx context.Context, orderId string, tradedAmount float64, turnover float64, status int) error
}
