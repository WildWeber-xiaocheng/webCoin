package dao

import (
	"context"
	"gorm.io/gorm"
	"market/internal/model"
	"webCoin-common/msdb"
	"webCoin-common/msdb/gorms"
)

type ExchangeCoinDao struct {
	conn *gorms.GormConn
}

func (d *ExchangeCoinDao) FindBySymbol(ctx context.Context, symbol string) (list *model.ExchangeCoin, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeCoin{}).Where("symbol=?", symbol).Take(&list).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (d *ExchangeCoinDao) FindVisible(ctx context.Context) (list []*model.ExchangeCoin, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.ExchangeCoin{}).Where("visible=?", 1).Find(&list).Error
	return
}

func NewExchangeCoinDao(db *msdb.MsDB) *ExchangeCoinDao {
	return &ExchangeCoinDao{
		conn: gorms.New(db.Conn),
	}
}
