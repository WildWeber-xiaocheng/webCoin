package domain

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"grpc-common/market/mclient"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
	"webCoin-common/msdb"
)

type MemberWalletDomain struct {
	memberWalletRepo repo.MemberWalletRepo
}

func (d *MemberWalletDomain) FindWalletBySymbol(ctx context.Context, id int64, name string, coin *mclient.Coin) (*model.MemberWalletCoin, error) {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, id, name)
	if err != nil {
		return nil, err
	}
	if mw == nil {
		//用户没有对应币种的钱包，则新建并存储
		mw, walletCoin := model.NewMemberWallet(id, coin)
		err := d.memberWalletRepo.Save(ctx, mw)
		if err != nil {
			return nil, err
		}
		return walletCoin, nil
	}
	nwc := &model.MemberWalletCoin{}
	copier.Copy(nwc, mw)
	nwc.Coin = coin
	return nwc, nil
}

// 冻结资金
func (d *MemberWalletDomain) Freeze(ctx context.Context, userId int64, money float64, symbol string) error {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	if mw.Balance < money {
		return errors.New("余额不足")
	}
	err = d.memberWalletRepo.UpdateFreeze(ctx, userId, money, symbol)
	return err
}

func NewMemberWalletDomain(db *msdb.MsDB) *MemberWalletDomain {
	return &MemberWalletDomain{
		dao.NewMemberWalletDao(db),
	}
}
