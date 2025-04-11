package dao

import (
	"context"
	"market/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KlineDao struct {
	db *mongo.Database
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}

func (d *KlineDao) FindBySymbol(ctx context.Context, symbol, period string, count int64) ([]*model.Kline, error) {
	
	collection := d.db.Collection((&model.Kline{}).Table(symbol, period))
	cursor,err := collection.Find(ctx,bson.D{{}},&options.FindOptions{
		Limit: &count,
		Sort: bson.D{{Key: "time", Value: -1}},
	})
	if err != nil {
		return nil, err
	}
	var klines []*model.Kline

	err = cursor.All(ctx, &klines)
	if err != nil {
		return nil, err
	}
	return klines, nil
}

func (d *KlineDao) FindBySymbolTime(ctx context.Context, symbol, period string, from, end int64, sort string) (list []*model.Kline, err error) {
	collection := d.db.Collection((&model.Kline{}).Table(symbol, period))
	sortInt := -1
	if "asc" == sort {
		sortInt = 1
	}
	cur, err := collection.Find(ctx,
		bson.D{{"time", bson.D{{"$gte", from}, {"$lte", end}}}},
		&options.FindOptions{
			Sort: bson.D{{"time", sortInt}},
		})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &list)
	if err != nil {
		return nil, err
	}
	return
}
