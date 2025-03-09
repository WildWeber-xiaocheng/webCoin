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

func (e *ExchangeOrderDao) Save(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder) error {
	e.conn = conn.(*gorms.GormConn)
	tx := e.conn.Tx(ctx)
	err := tx.Save(&order).Error
	return err
}

func (e *ExchangeOrderDao) FindCurrentTradingCount(ctx context.Context, userId int64, symbol string, direction int) (total int64, err error) {
	session := e.conn.Session(ctx)
	err = session.
		Model(&model.ExchangeOrder{}).
		Where("symbol = ? and member_id = ? and direction = ? and status = ?", symbol, userId, direction, model.Trading).
		Count(&total).Error
	return
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
