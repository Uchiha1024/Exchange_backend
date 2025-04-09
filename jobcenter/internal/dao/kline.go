package dao

import (
	"context"
	"jobcenter/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KlineDao struct {
	db *mongo.Database
}

func NewKlineDao(db *mongo.Database) *KlineDao {
	return &KlineDao{
		db: db,
	}
}

func (d *KlineDao) SaveBatch(ctx context.Context, data []*model.Kline, symbol, period string) error {
	 // 创建一个空的K线模型实例，用于获取集合名称
	mk := &model.Kline{}
	// 获取MongoDB集合，集合名称通过symbol和period动态生成
    // 例如：kline_btc_usdt_1min
	collection := d.db.Collection(mk.Table(symbol, period))
	// 创建一个interface{}切片，用于存储要插入的文档
    // 长度与输入数据切片相同
	ds := make([]interface{}, len(data))
	// 将K线数据转换为interface{}类型
    // MongoDB的InsertMany方法需要interface{}类型的切片
	for i, v := range data {
		ds[i] = v
	}
	// 使用InsertMany方法批量插入数据到MongoDB
    // 返回插入结果和可能的错误
	_, err := collection.InsertMany(ctx, ds)

	return err

}

func (d *KlineDao) DeleteGtTime(ctx context.Context, time int64, symbol string, period string) error {
	mk := &model.Kline{}
	collection := d.db.Collection(mk.Table(symbol, period))
	deleteResult, err := collection.DeleteMany(ctx, bson.D{{"time", bson.D{{"$gte", time}}}})
	if err != nil {
		return err
	}
	logx.Infof("%s %s 删除了%d条数据 \n", symbol, period, deleteResult.DeletedCount)

	return nil
}
