package logic

import (
	"context"
	"grpc-common/ucenter/types/asset"
	"mscoin-common/pages"
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



func (l *AssetLogic) ResetAddress(req *types.AssetReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	_, err := l.svcCtx.UCAssetRpc.ResetAddress(ctx, &asset.AssetReq{
		UserId:   userId,
		CoinName: req.Unit,
	})
	if err != nil {
		logx.Errorf("RPC-ResetAddress error: %v", err)
		return "", err
	}
	return "success", nil
}



func (l *AssetLogic) FindTransaction(req *types.AssetReq) (*pages.PageResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	resp, err := l.svcCtx.UCAssetRpc.FindTransaction(ctx, &asset.AssetReq{
		UserId:    userId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageNo:    int64(req.PageNo),
		PageSize:  int64(req.PageSize),
		Symbol:    req.Symbol,
		Type:      req.Type,
	})
	if err != nil {
		logx.Errorf("RPC-FindTransaction error: %v", err)
		return nil, err
	}
	total := resp.Total
	respList := make([] any, len(resp.List))
	for i, v := range resp.List {
		respList[i] = v
	}
	return pages.New(respList, int64(req.PageNo), int64(req.PageSize), total), nil
	
}
