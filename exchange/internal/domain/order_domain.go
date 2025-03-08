package domain

import (
	"context"
	"errors"
	"exchange/internal/repo"
	"webCoin-common/msdb"
)

type ExchangeOrderDomain struct {
	coinRepo repo.ExchangeOrderRepo
}

func (d *ExchangeOrderDomain) FindExchangeOrderInfo(ctx context.Context, unit string) (*model.ExchangeOrder, error) {
	coin, err := d.coinRepo.FindByUnit(ctx, unit)
	coin.ColdWalletAddress = ""
	if coin == nil {
		return nil, errors.New("币种不存在" + unit)
	}

	return coin, err
}

func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{coinRepo: dao.NewExchangeOrderDao(db)}
}
