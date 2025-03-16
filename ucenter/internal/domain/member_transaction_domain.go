package domain

import (
	"context"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
	"webCoin-common/msdb"
)

type MemberTransactionDomain struct {
	memberTransactionRepo repo.MemberTransactionRepo
}

func (d *MemberTransactionDomain) FindTransaction(
	ctx context.Context,
	pageNo int64,
	pageSize int64,
	userId int64,
	symbol string,
	startTime string,
	endTime string,
	t string) ([]*model.MemberTransactionVo, int64, error) {
	list, total, err := d.memberTransactionRepo.FindTransaction(ctx, int(pageNo), int(pageSize), userId, startTime, endTime, symbol, t)
	if err != nil {
		return nil, total, err
	}
	var voList = make([]*model.MemberTransactionVo, len(list))
	for i, v := range list {
		voList[i] = v.ToVo()
	}
	return voList, total, nil
}

func NewMemberTransactionDomain(db *msdb.MsDB) *MemberTransactionDomain {
	return &MemberTransactionDomain{
		memberTransactionRepo: dao.NewMemberTransactionDao(db),
	}
}
