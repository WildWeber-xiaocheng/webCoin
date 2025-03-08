package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"market/internal/model"
)

type KlineDao struct {
	db *mongo.Database
}

// 全部数据 按照时间降序
func (d *KlineDao) FindBySymbol(ctx context.Context, symbol, period string, count int64) (list []*model.Kline, err error) {
	mk := &model.Kline{}
	collection := d.db.Collection(mk.Table(symbol, period))
	//1是升序 -1 是降序
	cur, err := collection.Find(ctx, bson.D{{}}, &options.FindOptions{
		Limit: &count,
		Sort:  bson.D{{"time", -1}},
	})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &list)
	if err != nil {
		return
	}
	return
}

// 按照时间范围查询
// from: 起始时间 end: 结束时间
func (d *KlineDao) FindBySymbolTime(ctx context.Context, symbol, period string, from, end int64, sort string) (list []*model.Kline, err error) {
	mk := &model.Kline{}
	collection := d.db.Collection(mk.Table(symbol, period))
	//1是升序 -1 是降序
	sortInt := -1
	if "asc" == sort {
		sortInt = 1
	}
	cur, err := collection.Find(ctx, bson.D{{Key: "time", Value: bson.D{{"$gte", from}, {"$lte", end}}}}, &options.FindOptions{
		Sort: bson.D{{"time", sortInt}},
	})
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &list)
	if err != nil {
		return
	}
	return
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
