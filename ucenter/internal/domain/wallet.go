package domain

import (
	"context"
	"grpc-common/market/mclient"
	"mscoin-common/msdb"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MemberWalletDomain struct {
	memberWalletRepo repo.MemberWalletRepo
}

func NewMemberWalletDomain(db *msdb.MsDB) *MemberWalletDomain {
	return &MemberWalletDomain{
		memberWalletRepo: dao.NewMemberWalletDao(db),
	}
}

func (m *MemberWalletDomain) FindByIdAndCoinName(ctx context.Context, memId int64, coinName string, coin *mclient.Coin) (*model.MemberWalletCoin, error) {
	mw, err := m.memberWalletRepo.FindByIdAndCoinName(ctx, memId, coinName)
	if err != nil {
		logx.Errorf("DOMAIN-FindByIdAndCoinName - ERROR: %v", err)
		return nil, err
	}
	if mw == nil {
		// 新建并存储
		mw, walletCoint := model.NewMemberWallet(memId, coin)
		err = m.memberWalletRepo.Save(ctx, mw)
		if err != nil {
			logx.Errorf("DOMAIN-FindByIdAndCoinName - ERROR: %v", err)
			return nil, err
		}
		return walletCoint, nil
	}
	mwc := &model.MemberWalletCoin{}
	if err := copier.Copy(mwc, mw); err != nil {
		return nil, err // 或者适当的错误处理
	}
	mwc.Coin = coin
	return mwc, nil
}
