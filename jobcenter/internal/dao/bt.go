package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"jobcenter/internal/model"
)

type BtcTransactionDao struct {
	db *mongo.Database
}

func (b *BtcTransactionDao) FindByTxId(txId string) (bt *model.BitCoinTransaction, err error) {
	collection := b.db.Collection("bitcoin_transaction")
	err = collection.FindOne(
		context.Background(),
		bson.D{{"txId", txId}}).Decode(&bt)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return
}

func (b *BtcTransactionDao) Save(bt *model.BitCoinTransaction) error {
	collection := b.db.Collection("bitcoin_transaction")
	_, err := collection.InsertOne(
		context.Background(),
		&bt)
	return err
}

func NewBtcTransactionDao(db *mongo.Database) *BtcTransactionDao {
	return &BtcTransactionDao{
		db: db,
	}
}
