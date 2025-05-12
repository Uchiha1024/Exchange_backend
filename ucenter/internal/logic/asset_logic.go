package logic

import (
	"context"

	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"mscoin-common/bc"
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
	MemberTransactionDomain *domain.MemberTransactionDomain
}

func NewAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetLogic {
	return &AssetLogic{
		ctx:                ctx,
		svcCtx:             svcCtx,
		Logger:             logx.WithContext(ctx),
		CaptchaDomain:      domain.NewCaptchaDomain(),
		MemberDomain:       domain.NewMemberDomain(svcCtx.Db),
		memberWalletDomain: domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.Cache),
		MemberTransactionDomain: domain.NewMemberTransactionDomain(svcCtx.Db),
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

func (l *AssetLogic) FindWallet(in *asset.AssetReq) (*asset.MemberWalletList, error) {
	//根据用户id查询用户的钱包 循环钱包信息 根据币种 查询币种详情
	memberWallets, err := l.memberWalletDomain.FindWallet(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var list []*asset.MemberWallet
	err = copier.Copy(&list, memberWallets)
	return &asset.MemberWalletList{
		List: list,
	}, err

}

func (l *AssetLogic) ResetAddress(in *asset.AssetReq) (*asset.AssetResp, error) {
	//查询用户的钱包 检查address是否为空 如果未空 生成地址 进行更新
	memberWallet, err := l.memberWalletDomain.FindWalletByMemIdAndCoin(l.ctx, in.UserId, in.CoinName)
	if err != nil {
		return nil, err
	}
	if in.CoinName == "BTC" {
		if memberWallet.Address == "" {
			wallet, err := bc.NewWallet()
			if err != nil {
				return nil, err
			}
			address := wallet.GetTestAddress()
			priKey := wallet.GetPriKey()
			memberWallet.AddressPrivateKey = priKey
			memberWallet.Address = string(address)
			err = l.memberWalletDomain.UpdateAddress(l.ctx, memberWallet)
			if err != nil {
				return nil, err
			}
		}
	}
	return &asset.AssetResp{}, nil

}

func (l *AssetLogic) FindTransaction(in *asset.AssetReq) (*asset.MemberTransactionList, error) {
	//查询所有的充值记录 分页查询
	memberTransactionVos, total, err := l.MemberTransactionDomain.FindTransaction(
		l.ctx,
		in.UserId,
		in.PageNo,
		in.PageSize,
		in.Symbol,
		in.StartTime,
		in.EndTime,
		in.Type,
	)
	if err != nil {
		return nil, err
	}
	var list []*asset.MemberTransaction
	copier.Copy(&list, memberTransactionVos)
	return &asset.MemberTransactionList{
		List:  list,
		Total: total,
	}, nil
}
