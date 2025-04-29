package domain

import (
	"context"
	"errors"
	"grpc-common/market/mclient"
	"mscoin-common/msdb"
	"mscoin-common/msdb/tran"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type MemberWalletDomain struct {
	memberWalletRepo repo.MemberWalletRepo
	transaction      tran.Transaction
	marketRpc        mclient.Market
	redisCache       cache.Cache
}

func NewMemberWalletDomain(db *msdb.MsDB, marketRpc mclient.Market, redisCache cache.Cache) *MemberWalletDomain {
	return &MemberWalletDomain{
		memberWalletRepo: dao.NewMemberWalletDao(db),
		transaction:      tran.NewTransaction(db.Conn),
		marketRpc:        marketRpc,
		redisCache:       redisCache,
	}
}

func (m *MemberWalletDomain) Freeze(ctx context.Context, conn msdb.DbConn, userId int64, money float64, symbol string) error {

	// 获取钱包信息
	mw, err := m.memberWalletRepo.FindByIdAndCoinName(ctx, userId, symbol)
	if err != nil {
		logx.Errorf("DOMAIN-Freeze - ERROR: %v", err)
		return err
	}
	if mw.Balance < money {
		return errors.New("余额不足")
	}
	err = m.memberWalletRepo.UpdateFreeze(ctx, conn, userId, symbol, money)
	if err != nil {
		logx.Errorf("DOMAIN-Freeze - ERROR: %v", err)
		return err
	}
	return nil


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
