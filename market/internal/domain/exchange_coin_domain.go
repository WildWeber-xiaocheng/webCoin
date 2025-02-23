package domain

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
	"webCoin-common/msdb"
)

type ExchangeCoinDomain struct {
	exchangeCoinRepo repo.ExchangeCoinRepo //repo小写，是防止外部直接操作repo，而不用domain
}

func NewExchangeCoinDomain(db *msdb.MsDB) *ExchangeCoinDomain {
	return &ExchangeCoinDomain{
		exchangeCoinRepo: dao.NewExchangeCoinDao(db),
	}
}

func (d *ExchangeCoinDomain) FindVisible(ctx context.Context) []*model.ExchangeCoin {
	list, err := d.exchangeCoinRepo.FindVisible(ctx)
	if err != nil {
		logx.Error(err)
		return []*model.ExchangeCoin{}
	}
	return list
}
