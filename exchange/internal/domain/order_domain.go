package domain

import (
	"context"
	"exchange/internal/dao"
	"exchange/internal/model"
	"exchange/internal/repo"
	"webCoin-common/msdb"
)

type ExchangeOrderDomain struct {
	orderRepo repo.ExchangeOrderRepo
}

func (d *ExchangeOrderDomain) FindOrderHistory(
	ctx context.Context,
	symbol string,
	page int64,
	size int64,
	memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderHistory(ctx, symbol, page, size, memberId)
	if err != nil {
		return nil, 0, err
	}
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func (d *ExchangeOrderDomain) FindOrderCurrent(
	ctx context.Context,
	symbol string,
	page int64,
	size int64,
	memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderCurrent(ctx, symbol, page, size, memberId)
	if err != nil {
		return nil, 0, err
	}
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{orderRepo: dao.NewExchangeOrderDao(db)}
}
