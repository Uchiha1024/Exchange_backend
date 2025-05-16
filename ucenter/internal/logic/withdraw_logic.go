package logic

import (
	"context"
	"grpc-common/ucenter/types/withdraw"
	"mscoin-common/msdb/tran"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	memberAddressDomain *domain.MemberAddressDomain
	memberDomain        *domain.MemberDomain
	memberWalletDomain  *domain.MemberWalletDomain
	transaction         tran.Transaction
	// withdrawDomain      *domain.WithdrawDomain
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		memberAddressDomain: domain.NewMemberAddressDomain(svcCtx.Db),
		memberDomain:        domain.NewMemberDomain(svcCtx.Db),
		transaction:         tran.NewTransaction(svcCtx.Db.Conn),
		memberWalletDomain:  domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.Cache),
		// withdrawDomain:      domain.NewWithdrawDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.BitcoinAddress),
	}
}

func (l *WithdrawLogic) FindAddressByCoinId(req *withdraw.WithdrawReq) (*withdraw.AddressSimpleList, error) {
	list, err := l.memberAddressDomain.FindAddressList(l.ctx, req.UserId, req.CoinId)
	if err != nil {
		return nil, err
	}
	var addressList []*withdraw.AddressSimple
	copier.Copy(&addressList, list)
	return &withdraw.AddressSimpleList{
		List: addressList,
	}, nil
}
