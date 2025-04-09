package repo

import (
	"context"
	"ucenter/model"
)

type MemberRepo interface {
	FindByPhone(ctx context.Context, phone string) (*model.Member, error)
	Save(ctx context.Context, mem *model.Member) error
	UpdateLoginCount(ctx context.Context, id int64, step int) error

}
