package logic

import (
	"context"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/withdraw"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type Withdraw struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Withdraw {
	return &Withdraw{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (w *Withdraw) QueryWithdrawCoin(req *types.WithdrawReq) ([]*types.WithdrawWalletInfo, error) {

	userId := w.ctx.Value("userId").(int64)
	// 查询 币种 信息
	coinList, err := w.svcCtx.MarketRpc.FindAllCoin(w.ctx, &market.MarketReq{})
	if err != nil {
		return nil, err
	}
	// 创建map 存储 币种 信息
	coinMap := make(map[string]*market.Coin)
	for _, coin := range coinList.List {
		coinMap[coin.Unit] = coin
	}

	//2. 根据用户id 查询用户的钱包信息
	walletList, err := w.svcCtx.UCAssetRpc.FindWallet(w.ctx, &asset.AssetReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	//3. 组装数据[]*types.WithdrawWalletInfo
	wwList := make([]*types.WithdrawWalletInfo, len(walletList.List))
	for i, wallet := range walletList.List {
		var ww types.WithdrawWalletInfo
		coin := coinMap[wallet.Coin.Unit]
		ww.Balance = wallet.Balance
		ww.WithdrawScale = int(coin.WithdrawScale)
		ww.MaxTxFee = coin.MaxTxFee
		ww.MinTxFee = coin.MinTxFee
		ww.MaxAmount = coin.MaxWithdrawAmount
		ww.MinAmount = coin.MinWithdrawAmount
		ww.Name = coin.GetName()
		ww.NameCn = coin.NameCn
		ww.Threshold = coin.WithdrawThreshold
		ww.Unit = coin.Unit
		ww.AccountType = int(coin.AccountType)
		if coin.CanAutoWithdraw == 0 {
			ww.CanAutoWithdraw = "true"
		} else {
			ww.CanAutoWithdraw = "false"
		}
		//提币地址的赋值
		addressSimpleList, err := w.svcCtx.UCWithdrawRpc.FindAddressByCoinId(w.ctx, &withdraw.WithdrawReq{
			UserId: userId,
			CoinId: int64(coin.Id),
		})
		if err != nil {
			return nil, err
		}
		var addressList []types.AddressSimple
		copier.Copy(&addressList, addressSimpleList.List)
		ww.Addresses = addressList
		wwList[i] = &ww
	}
	return wwList, nil

}
