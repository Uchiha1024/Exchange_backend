package domain

import (
	"context"
	"mscoin-common/msdb"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberAddressDomain struct {
	memberAddressRepo repo.MemberAddressRepo
}

func NewMemberAddressDomain(db *msdb.MsDB) *MemberAddressDomain {
	return &MemberAddressDomain{
		memberAddressRepo: dao.NewMemberAddressDao(db),
	}
}

func (d *MemberAddressDomain) FindAddressList(
	ctx context.Context,
	userId int64,
	coinId int64) ([]*model.MemberAddress, error) {
	return d.memberAddressRepo.FindByMemIdAndCoinId(ctx, userId, coinId)
}
