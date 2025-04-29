package logic

import (
	"context"
	"errors"
	"exchange/internal/domain"
	"exchange/internal/model"
	"exchange/internal/svc"
	"fmt"
	"grpc-common/exchange/types/order"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/member"
	"mscoin-common/msdb"
	"mscoin-common/msdb/tran"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeOrderDomain *domain.ExchangeOrderDomain
	transaction         tran.Transaction
	kafkaDomain         *domain.KafkaDomain
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	orderDomain := domain.NewExchangeOrderDomain(svcCtx.Db)
	return &ExchangeOrderLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		exchangeOrderDomain: orderDomain,
		transaction:         tran.NewTransaction(svcCtx.Db.Conn),
		kafkaDomain:         domain.NewKafkaDomain(svcCtx.KafkaClient, orderDomain),
	}
}

func (l *ExchangeOrderLogic) FindOrderHistory(req *order.OrderReq) (*order.OrderRes, error) {
	voList, total, err := l.exchangeOrderDomain.FindOrderHistory(l.ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory", logx.Field("error", err))
		return nil, err
	}
	var list []*order.ExchangeOrder
	err = copier.Copy(&list, &voList)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory Copier Error", logx.Field("error", err))
		return nil, err
	}
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) FindOrderCurrent(req *order.OrderReq) (*order.OrderRes, error) {
	voList, total, err := l.exchangeOrderDomain.FindOrderCurrent(l.ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory", logx.Field("error", err))
		return nil, err
	}
	var list []*order.ExchangeOrder
	err = copier.Copy(&list, &voList)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory Copier Error", logx.Field("error", err))
		return nil, err
	}
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) AddOrder(req *order.OrderReq) (*order.AddOrderRes, error) {
	//添加订单 发布委托
	//1.首先检查参数是否合法
	memberRes, err := l.svcCtx.MemberRpc.FindMemberById(l.ctx, &member.MemberReq{
		MemberId: req.UserId,
	})
	if err != nil {
		logx.Errorw("MemberRpc-FindMemberById-ERROR", logx.Field("err", err))
		return nil, err
	}
	if memberRes.TransactionStatus == 0 {
		return nil, errors.New("此用户已经被禁止交易")
	}
	if req.Type == model.TypeMap[model.LimitPrice] && req.Price <= 0 {
		return nil, errors.New("限价模式下价格不能小于等于0")
	}
	if req.Amount <= 0 {
		return nil, errors.New("数量不能小于等于0")
	}

	exchangeCoin, err := l.svcCtx.MarketRpc.FindSymbolInfo(l.ctx, &market.MarketReq{
		Symbol: req.Symbol,
	})
	if err != nil {
		logx.Errorw("MarketRpc-FindSymbolInfo-ERROR", logx.Field("err", err))
		return nil, errors.New("nonsupport coin")
	}
	if exchangeCoin.Exchangeable != 1 && exchangeCoin.Enable != 1 {
		return nil, errors.New("coin forbidden")
	}
	//基准币 BTC/USDT
	baseSymbol := exchangeCoin.GetBaseSymbol()
	//交易币
	coinSymbol := exchangeCoin.GetCoinSymbol()
	cc := baseSymbol
	if req.Direction == model.DirectionMap[model.SELL] {
		//根据交易币查询
		cc = coinSymbol
	}
	coinInfo, err := l.svcCtx.MarketRpc.FindCoinInfo(l.ctx, &market.MarketReq{
		Unit: cc,
	})
	if err != nil || coinInfo == nil {
		logx.Errorw("MarketRpc-FindCoinInfo-ERROR", logx.Field("err", err))
		return nil, errors.New("nonsupport coin")
	}
	if req.Type == model.TypeMap[model.MarketPrice] && req.Direction == model.DirectionMap[model.BUY] {
		if exchangeCoin.GetMinTurnover() > 0 && req.Amount < float64(exchangeCoin.GetMinTurnover()) {
			return nil, errors.New("成交额至少是" + fmt.Sprintf("%d", exchangeCoin.GetMinTurnover()))
		}
	} else {
		if exchangeCoin.GetMaxVolume() > 0 && exchangeCoin.GetMaxVolume() < req.Amount {
			return nil, errors.New("数量超出" + fmt.Sprintf("%f", exchangeCoin.GetMaxVolume()))
		}
		if exchangeCoin.GetMinVolume() > 0 && exchangeCoin.GetMinVolume() > req.Amount {
			return nil, errors.New("数量不能低于" + fmt.Sprintf("%f", exchangeCoin.GetMinVolume()))
		}
	}

	//查询用户钱包 BTC/USDT
	baseWallet, err := l.svcCtx.AssetRpc.FindWalletBySymbol(l.ctx, &asset.AssetReq{
		UserId:   req.UserId,
		CoinName: baseSymbol,
	})
	if err != nil {
		return nil, errors.New("no wallet")
	}
	exCoinWallet, err := l.svcCtx.AssetRpc.FindWalletBySymbol(l.ctx, &asset.AssetReq{
		UserId:   req.UserId,
		CoinName: coinSymbol,
	})
	if err != nil {
		return nil, errors.New("no wallet")
	}
	if baseWallet.IsLock == 1 || exCoinWallet.IsLock == 1 {
		return nil, errors.New("wallet locked")
	}
	if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.GetMinSellPrice() > 0 {
		if req.Price < exchangeCoin.GetMinSellPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, errors.New("不能低于最低限价:" + fmt.Sprintf("%f", exchangeCoin.GetMinSellPrice()))
		}
	}
	if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.GetMaxBuyPrice() > 0 {
		if req.Price > exchangeCoin.GetMaxBuyPrice() || req.Type == model.TypeMap[model.MarketPrice] {
			return nil, errors.New("不能低于最高限价:" + fmt.Sprintf("%f", exchangeCoin.GetMaxBuyPrice()))
		}
	}
	//是否启用了市价买卖
	if req.Type == model.TypeMap[model.MarketPrice] {
		if req.Direction == model.DirectionMap[model.BUY] && exchangeCoin.EnableMarketBuy == 0 {
			return nil, errors.New("不支持市价购买")
		} else if req.Direction == model.DirectionMap[model.SELL] && exchangeCoin.EnableMarketSell == 0 {
			return nil, errors.New("不支持市价出售")
		}
	}
	//限制委托数量
	count, err := l.exchangeOrderDomain.FindCurrentTradingCount(l.ctx, req.UserId, req.Symbol, req.Direction)
	if err != nil {
		return nil, err
	}
	if exchangeCoin.GetMaxTradingOrder() > 0 && count >= exchangeCoin.GetMaxTradingOrder() {
		return nil, errors.New("超过最大挂单数量 " + fmt.Sprintf("%d", exchangeCoin.GetMaxTradingOrder()))
	}
	// 生成订单
	exchangeOrder := model.NewOrder()
	exchangeOrder.MemberId = req.UserId
	exchangeOrder.Symbol = req.Symbol
	exchangeOrder.BaseSymbol = baseSymbol
	exchangeOrder.CoinSymbol = coinSymbol
	typeCode := model.TypeMap.Code(req.Type)
	exchangeOrder.Type = typeCode
	directionCode := model.DirectionMap.Code(req.Direction)
	exchangeOrder.Direction = directionCode
	if exchangeOrder.Type == model.MarketPrice {
		exchangeOrder.Price = 0

	} else {
		exchangeOrder.Price = req.Price
	}
	exchangeOrder.UseDiscount = "0"
	exchangeOrder.Amount = req.Amount
	//保存订单到数据库，发送消息到kafka，ucenter 钱包服务 接收到消息 进行资金的冻结
	//AddOrder 保存订单 计算所需要的钱
	err = l.transaction.Action(func(conn msdb.DbConn) error {
		money, err := l.exchangeOrderDomain.AddOrder(l.ctx, conn, exchangeOrder, exchangeCoin, baseWallet, exCoinWallet)
		if err != nil {
			return errors.New("订单提交失败")
		}
		// 发送消息到kafka
		err = l.kafkaDomain.SendOrder("add-exchange-order",
			req.UserId,
			exchangeOrder.OrderId,
			money,
			req.Symbol,
			exchangeOrder.Direction,
			baseSymbol,
			coinSymbol)
		if err != nil {
			return errors.New("发消息失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &order.AddOrderRes{
		OrderId: exchangeOrder.OrderId,
	}, nil
}

func (l *ExchangeOrderLogic) FindByOrderId(req *order.OrderReq) (*order.ExchangeOrderOrigin, error) {
	orderResp, err := l.exchangeOrderDomain.FindByOrderId(l.ctx, req.OrderId)
	if err != nil {
		logx.Errorw("Logic-FindByOrderId", logx.Field("error", err))
		return nil, err
	}
	resp := &order.ExchangeOrderOrigin{}
	err = copier.Copy(resp, orderResp)
	if err != nil {
		logx.Errorw("Logic-FindByOrderId-Copier Error", logx.Field("error", err))
		return nil, err
	}
	return resp, nil

}

func (l *ExchangeOrderLogic) CancelOrder(req *order.OrderReq) (*order.CancelOrderRes, error) {
	l.exchangeOrderDomain.UpdateStatusCancel(l.ctx, req.OrderId)
	return nil, nil

}
