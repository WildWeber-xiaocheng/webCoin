package dao

import (
	"context"
	"exchange/internal/model"
	"webCoin-common/msdb"
	"webCoin-common/msdb/gorms"
)

type ExchangeOrderDao struct {
	conn *gorms.GormConn
}

func (e *ExchangeOrderDao) FindOrderHistory(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := e.conn.Session(ctx)
	index := (page - 1) * size
	err = session.
		Model(&model.ExchangeOrder{}).
		Where("symbol = ? and member_id = ?", symbol, memberId).
		Limit(int(size)).
		Offset(int(index)).
		Find(&list).Error
	err = session.
		Model(&model.ExchangeOrder{}).
		Where("symbol = ? and member_id = ?", symbol, memberId).
		Count(&total).Error
	return
}

func (e *ExchangeOrderDao) FindOrderCurrent(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := e.conn.Session(ctx)
	index := (page - 1) * size
	err = session.
		Model(&model.ExchangeOrder{}).
		Where("symbol = ? and member_id = ? and status = ?", symbol, memberId, model.Trading).
		Limit(int(size)).
		Offset(int(index)).
		Find(&list).Error
	err = session.
		Model(&model.ExchangeOrder{}).
		Where("symbol = ? and member_id = ? and status = ?", symbol, memberId, model.Trading).
		Count(&total).Error
	return
}

func NewExchangeOrderDao(db *msdb.MsDB) *ExchangeOrderDao {
	return &ExchangeOrderDao{
		conn: gorms.New(db.Conn),
	}
}
