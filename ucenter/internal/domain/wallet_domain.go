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
	"webCoin-common/msdb/tran"
)

type MemberWalletDomain struct {
	memberWalletRepo repo.MemberWalletRepo
	transaction      tran.Transaction
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
func (d *MemberWalletDomain) Freeze(ctx context.Context, conn msdb.DbConn, userId int64, money float64, symbol string) error {
	mw, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, userId, symbol)
	if err != nil {
		return err
	}
	if mw.Balance < money {
		return errors.New("余额不足")
	}
	//只在更新数据的时候使用事务
	err = d.memberWalletRepo.UpdateFreeze(ctx, conn, userId, money, symbol)
	return err
}

func (d *MemberWalletDomain) FindWalletByMemIdAndCoin(ctx context.Context, memberId int64, coinName string) (*model.MemberWallet, error) {
	memberWallet, err := d.memberWalletRepo.FindByIdAndCoinName(ctx, memberId, coinName)
	if err != nil {
		return nil, err
	}
	return memberWallet, nil
}

func (d *MemberWalletDomain) UpdateWalletCoinAndBase(ctx context.Context, baseWallet *model.MemberWallet, coinWallet *model.MemberWallet) error {
	return d.transaction.Action(func(conn msdb.DbConn) error {
		err := d.memberWalletRepo.UpdateWallet(ctx, conn, baseWallet.Id, baseWallet.Balance, baseWallet.FrozenBalance)
		if err != nil {
			return err
		}
		err = d.memberWalletRepo.UpdateWallet(ctx, conn, coinWallet.Id, coinWallet.Balance, coinWallet.FrozenBalance)
		if err != nil {
			return err
		}
		return nil
	})
}

func NewMemberWalletDomain(db *msdb.MsDB) *MemberWalletDomain {
	return &MemberWalletDomain{
		memberWalletRepo: dao.NewMemberWalletDao(db),
		transaction:      tran.NewTransaction(db.Conn),
	}
}
