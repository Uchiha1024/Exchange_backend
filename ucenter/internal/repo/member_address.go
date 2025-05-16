package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberAddressRepo interface {
	FindByMemIdAndCoinId(ctx context.Context, memId int64, coinId int64) ([]*model.MemberAddress, error)
}
