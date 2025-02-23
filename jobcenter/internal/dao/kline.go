package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"jobcenter/internal/model"
	"log"
)

type KlineDao struct {
	db *mongo.Database
}

// 批量删除
func (d *KlineDao) DeleteGtTime(ctx context.Context, time int64, symbol string, period string) error {
	mk := &model.Kline{}
	collection := d.db.Collection(mk.Table(symbol, period))
	deleteResult, err := collection.DeleteMany(ctx, bson.D{{"time", bson.D{{"$gte", time}}}})
	if err != nil {
		return err
	}
	log.Printf("删除表%s，数量为：%d \n", "exchange_kline_"+symbol+"_"+period, deleteResult.DeletedCount)
	return nil
}

// 批量保存
func (d *KlineDao) SaveBatch(ctx context.Context, data []*model.Kline, symbol, period string) error {
	mk := &model.Kline{}
	collection := d.db.Collection(mk.Table(symbol, period))
	ds := make([]interface{}, len(data))
	for i, v := range data {
		ds[i] = v
	}
	//因为InsertMany需要传入[]interface{},所以这里转一下
	_, err := collection.InsertMany(ctx, ds)
	return err
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}
