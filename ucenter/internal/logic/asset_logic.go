package logic

import (
	"context"

	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AssetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	CaptchaDomain      *domain.CaptchaDomain
	MemberDomain       *domain.MemberDomain
	memberWalletDomain *domain.MemberWalletDomain
}

func NewAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetLogic {
	return &AssetLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		CaptchaDomain:      domain.NewCaptchaDomain(),
		MemberDomain:       domain.NewMemberDomain(svcCtx.Db),
		memberWalletDomain:      domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.Cache),

	}
}

func (l *AssetLogic) FindWalletBySymbol(in *asset.AssetReq) (*asset.MemberWallet, error) {
	//通过market rpc 进行coin表的查询 coin信息
	//通过钱包 查询对应币的钱包信息  coin_id  user_id 查询用户的钱包信息 组装信息
	coinInfo, err := l.svcCtx.MarketRpc.FindCoinInfo(l.ctx, &market.MarketReq{
		Unit: in.CoinName,
	})
	if err != nil {
		logx.Errorf("RPC-FindCoinInfo error: %v", err)
		return nil, err
	}
	memberWalletCoin, err := l.memberWalletDomain.FindByIdAndCoinName(l.ctx, in.UserId, in.CoinName, coinInfo)
	if err != nil {
		return nil, err
	}
	resp := &asset.MemberWallet{}
	err = copier.Copy(resp, memberWalletCoin)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
