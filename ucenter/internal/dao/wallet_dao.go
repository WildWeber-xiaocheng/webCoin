package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"ucenter/internal/model"
	"webCoin-common/msdb"
	"webCoin-common/msdb/gorms"
)

type MemberWalletDao struct {
	conn *gorms.GormConn
}

func (m *MemberWalletDao) UpdateFreeze(ctx context.Context, conn msdb.DbConn, userId int64, money float64, symbol string) error {
	m.conn = conn.(*gorms.GormConn)
	session := m.conn.Tx(ctx)
	query := "update member_wallet set balance=balance-?,frozen_balance=frozen_balance+? where member_id=? and coin_name=? and balance > ?"
	exec := session.Exec(query, money, money, userId, symbol, money)
	err := exec.Error
	if err != nil {
		return err
	}
	affected := exec.RowsAffected
	if affected <= 0 {
		return errors.New("no update row")
	}
	return nil
}

func (m *MemberWalletDao) Save(ctx context.Context, mw *model.MemberWallet) error {
	session := m.conn.Session(ctx)
	err := session.Save(&mw).Error
	return err
}

func (m *MemberWalletDao) FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).
		Where("member_id=? and coin_name=?", memId, coinName).
		Take(&mw).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func NewMemberWalletDao(db *msdb.MsDB) *MemberWalletDao {
	return &MemberWalletDao{
		conn: gorms.New(db.Conn),
	}
}
