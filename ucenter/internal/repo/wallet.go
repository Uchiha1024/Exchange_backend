package repo

import (
	"context"
	"mscoin-common/msdb"
	"ucenter/internal/model"
)

type MemberWalletRepo interface {
	Save(ctx context.Context, mw *model.MemberWallet) error
	FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error)
	UpdateFreeze(ctx context.Context, conn msdb.DbConn, memberId int64, symbol string, money float64) error
	UpdateWallet(ctx context.Context, conn msdb.DbConn, id int64, balance float64, frozenBalance float64) error
	FindByMemberId(ctx context.Context, memId int64) ([]*model.MemberWallet, error)
	UpdateAddress(ctx context.Context, wallet *model.MemberWallet) error
	FindAllAddress(ctx context.Context, name string) ([]string, error)
	FindByAddress(ctx context.Context, address string) (*model.MemberWallet, error)
	FindByIdAndCoinId(ctx context.Context, memberId int64, coinId int64) (*model.MemberWallet, error)
}
