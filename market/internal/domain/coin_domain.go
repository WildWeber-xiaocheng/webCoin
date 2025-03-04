package domain

import (
	"context"
	"errors"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
	"webCoin-common/msdb"
)

type CoinDomain struct {
	coinRepo repo.CoinRepo
}

func (d *CoinDomain) FindCoinInfo(ctx context.Context, unit string) (*model.Coin, error) {
	coin, err := d.coinRepo.FindByUnit(ctx, unit)
	coin.ColdWalletAddress = ""
	if coin == nil {
		return nil, errors.New("币种不存在" + unit)
	}

	return coin, err
}

func NewCoinDomain(db *msdb.MsDB) *CoinDomain {
	return &CoinDomain{coinRepo: dao.NewCoinDao(db)}
}
