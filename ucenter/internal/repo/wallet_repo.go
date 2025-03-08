package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberWalletRepo interface {
	Save(ctx context.Context, mw *model.MemberWallet) error
	FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error)
}
