package repo

import (
	"context"
	"ucenter/internal/model"
)

type MemberRepo interface {
	FindByPhone(ctx context.Context, phone string) (*model.Member, error)
	Save(ctx context.Context, mem *model.Member) error
	UpdateLoginCount(ctx context.Context, id int64, step int) error
	FindMemberById(ctx context.Context, memberId int64) (*model.Member, error)

}
