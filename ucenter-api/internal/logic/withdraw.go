package logic

import (
	"context"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/member"
	"grpc-common/ucenter/types/withdraw"
	"mscoin-common/pages"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	userId := w.ctx.Value("userId").(int64)
	// 查询 币种 信息
	coinList, err := w.svcCtx.MarketRpc.FindAllCoin(ctx, &market.MarketReq{})
	if err != nil {
		return nil, err
	}
	// 创建map 存储 币种 信息
	coinMap := make(map[string]*market.Coin)
	for _, coin := range coinList.List {
		coinMap[coin.Unit] = coin
	}

	//2. 根据用户id 查询用户的钱包信息
	walletList, err := w.svcCtx.UCAssetRpc.FindWallet(ctx, &asset.AssetReq{
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
		addressSimpleList, err := w.svcCtx.UCWithdrawRpc.FindAddressByCoinId(ctx, &withdraw.WithdrawReq{
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

func (w *Withdraw) SendCode(req *types.WithdrawReq) (string, error) {

	userId := w.ctx.Value("userId").(int64)
	// 查询用户信息,获取手机号
	userInfo, err := w.svcCtx.UCMemberRpc.FindMemberById(w.ctx, &member.MemberReq{
		MemberId: userId,
	})
	if err != nil {
		return "", err
	}
	phone := userInfo.MobilePhone
	// 根据手机号 发送验证码
	_, err = w.svcCtx.UCWithdrawRpc.SendCode(w.ctx, &withdraw.WithdrawReq{
		Phone: phone,
	})
	if err != nil {
		return "", err
	}
	return "success", nil

}

func (w *Withdraw) WithdrawCode(req *types.WithdrawReq) (string, error) {
	//1. 将参数传递给rpc服务 进行提现处理
	value := w.ctx.Value("userId").(int64)
	_, err := w.svcCtx.UCWithdrawRpc.WithdrawCode(w.ctx, &withdraw.WithdrawReq{
		UserId:     value,
		Unit:       req.Unit,
		JyPassword: req.JyPassword,
		Code:       req.Code,
		Address:    req.Address,
		Amount:     req.Amount,
		Fee:        req.Fee,
	})
	if err != nil {
		return "fail", err
	}
	return "success", nil
}

func (w *Withdraw) Record(req *types.WithdrawReq) (*pages.PageResult, error) {
	value := w.ctx.Value("userId").(int64)
	records, err := w.svcCtx.UCWithdrawRpc.WithdrawRecord(w.ctx, &withdraw.WithdrawReq{
		UserId:   value,
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := records.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, int64(req.Page), int64(req.PageSize), records.Total), nil
}
