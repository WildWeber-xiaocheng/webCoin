package repo

import (
	"context"
	"ucenter/internal/model"
	"webCoin-common/msdb"
)

type MemberWalletRepo interface {
	Save(ctx context.Context, mw *model.MemberWallet) error
	FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error)
	UpdateFreeze(ctx context.Context, conn msdb.DbConn, userId int64, money float64, symbol string) error
}
