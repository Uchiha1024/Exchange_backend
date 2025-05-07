package domain

import (
	"context"
	"errors"
	"exchange/internal/dao"
	"exchange/internal/model"
	"exchange/internal/repo"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/ucclient"
	"mscoin-common/msdb"
	"mscoin-common/op"
	"mscoin-common/tools"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderDomain struct {
	orderRepo repo.ExchangeOrderRepo
}



func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{
		orderRepo: dao.NewExchangeOrderDao(db),
	}
}

func (d *ExchangeOrderDomain) FindOrderHistory(ctx context.Context, symbol string, page int64, size int64, memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderHistory(ctx, symbol, page, size, memberId)
	if err != nil {
		logx.Errorw("Domain-FindOrderHistory", logx.Field("error", err))
		return nil, 0, err
	}
	voList := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		voList[i] = v.ToVo()

	}
	return voList, total, nil
}

func (d *ExchangeOrderDomain) FindOrderCurrent(ctx context.Context, symbol string, page int64, size int64, memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderCurrent(ctx, symbol, page, size, memberId)
	if err != nil {
		logx.Errorw("Domain-FindOrderCurrent", logx.Field("error", err))
	}
	voList := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		voList[i] = v.ToVo()

	}
	return voList, total, nil
}

func (d *ExchangeOrderDomain) FindCurrentTradingCount(ctx context.Context, userId int64, symbol string, direction string) (int64, error) {
	return d.orderRepo.FindCurrentTradingCount(ctx, userId, symbol, model.DirectionMap.Code(direction))
}

func (d *ExchangeOrderDomain) FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error) {
	order, err := d.orderRepo.FindOrderByOrderId(ctx, orderId)
	if err == nil && order == nil {
		return nil, errors.New("orderId:" + orderId + "订单不存在")
	}
	return order, nil
}

func (d *ExchangeOrderDomain) UpdateStatusCancel(ctx context.Context, orderId string) error {
	return d.orderRepo.UpdateStatusCancel(ctx, orderId)
}

func (d *ExchangeOrderDomain) UpdateOrderStatusTrading(ctx context.Context, orderId string) error {
	return d.orderRepo.UpdateOrderStatusTrading(ctx, orderId)
}


func (d *ExchangeOrderDomain) FindOrderListBySymbol(ctx context.Context, symbol string, status int) ([]*model.ExchangeOrder, error) {
	return d.orderRepo.FindOrderListBySymbol(ctx, symbol, status)
}

func (d *ExchangeOrderDomain) UpdateOrderComplete(context context.Context, orderInfo *model.ExchangeOrder) any {
	return d.orderRepo.UpdateOrderComplete(context, orderInfo.OrderId, orderInfo.TradedAmount, orderInfo.Turnover, orderInfo.Status)
}


func (d *ExchangeOrderDomain) AddOrder(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder, coin *mclient.ExchangeCoin,
	baseWallet *ucclient.MemberWallet,
	coinWallet *ucclient.MemberWallet) (float64, error) {
	order.Status = model.Init
	order.TradedAmount = 0
	order.Time = time.Now().UnixMilli()
	order.OrderId = tools.Unq("E")
	//交易的时候  coin.Fee 费率 手续费 我们做的时候 先不考虑手续费
	//买 花USDT 市价 price 0 冻结的直接就是amount  卖 BTC
	var money float64
	if order.Direction == model.BUY {
		if order.Type == model.MarketPrice {
			money = order.Amount
		} else {
			money = op.MulFloor(order.Price, order.Amount, 8)
		}
		if baseWallet.Balance < money {
			return 0, errors.New("余额不足")
		}
	} else {
		money = order.Amount
		if coinWallet.Balance < money {
			return 0, errors.New("余额不足")
		}
	}

	// 保存订单信息
	err := d.orderRepo.Save(ctx, conn, order)
	return money, err

}
