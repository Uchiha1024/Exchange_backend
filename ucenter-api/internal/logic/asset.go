package logic

import (
	"context"
	"grpc-common/ucenter/types/asset"
	"time"

	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetLogic {
	return &AssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetLogic) FindWalletBySymbol(req *types.AssetReq) (*types.MemberWallet, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	memberWallet, err := l.svcCtx.UCAssetRpc.FindWalletBySymbol(ctx, &asset.AssetReq{
		CoinName: req.CoinName,
		UserId:   userId,
	})

	if err != nil {
		logx.Errorf("RPC-FindWalletBySymbol error: %v", err)
		return nil, err
	}
	resp := &types.MemberWallet{}
	logx.Info("FindWalletBySymbol---memberWallet", memberWallet)
	if err := copier.Copy(resp, memberWallet); err != nil {
		return nil, err
	}

	return resp, nil

}

func (l *AssetLogic) FindWallet(req *types.AssetReq) ([]*types.MemberWallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	memberWallet, err := l.svcCtx.UCAssetRpc.FindWallet(ctx, &asset.AssetReq{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("RPC-FindWallet error: %v", err)
		return nil, err
	}
	var resp []*types.MemberWallet
	if err := copier.Copy(&resp, memberWallet.List); err != nil {
		return nil, err
	}
	return resp, nil
}
