package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberTransactionRepo interface {
	FindTransaction(
		ctx context.Context,
		pageNo int,
		pageSize int,
		memberId int64,
		startTime string,
		endTime string,
		symbol string,
		transactionType string) (list []*model.MemberTransaction, total int64, err error)
	FindByAmountAndTime(ctx context.Context, address string, value float64, time int64) (*model.MemberTransaction, error)
	Save(ctx context.Context, transaction *model.MemberTransaction) error
}
