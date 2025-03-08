package repo

import (
	"context"
	"exchange/internal/model"
)

type ExchangeOrderRepo interface {
	FindOrderPage(ctx context.Context, memberId int64) (coin *model.ExchangeOrder, err error)
}
