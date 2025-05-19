package repo

import (
	"context"
	"ucenter/internal/model"
)

type WithdrawRecordRepo interface {
	Save(ctx context.Context, record *model.WithdrawRecord) error
	UpdateSuccess(ctx context.Context, txId model.WithdrawRecord) error
	FindByUserId(ctx context.Context, userId int64, page int64, pageSize int64) ([]*model.WithdrawRecord, int64, error)
}