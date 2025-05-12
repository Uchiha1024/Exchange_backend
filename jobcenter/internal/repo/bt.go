package repo

import "jobcenter/internal/model"

type BtcTransactionRepo interface {
	FindByTxId(txId string) (*model.BitCoinTransaction, error)
	Save(bt *model.BitCoinTransaction) error
}
