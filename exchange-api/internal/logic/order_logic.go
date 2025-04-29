package logic

import (
	"context"
	"errors"
	"exchange-api/internal/svc"
	"exchange-api/internal/types"
	"grpc-common/exchange/types/order"
	"mscoin-common/pages"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderLogic) History(req *types.ExchangeReq) (*pages.PageResult, error) {
	ctx, cancel := context.WithTimeout(l.ctx, 10*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	orderResp, err := l.svcCtx.OrderRpc.FindOrderHistory(ctx, &order.OrderReq{
		Symbol:   req.Symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   userId,
	})
	if err != nil {
		logx.Errorw("OrderRpc-FindOrderHistory-ERROR", logx.Field("err", err))
		return nil, err
	}

	list := orderResp.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, req.PageNo, req.PageSize, orderResp.Total), nil

}

func (l *OrderLogic) Current(req *types.ExchangeReq) (*pages.PageResult, error) {
	ctx, cancel := context.WithTimeout(l.ctx, 10*time.Second)
	defer cancel()
	userId := l.ctx.Value("userId").(int64)
	symbol := req.Symbol
	orderRes, err := l.svcCtx.OrderRpc.FindOrderCurrent(ctx, &order.OrderReq{
		Symbol:   symbol,
		Page:     req.PageNo,
		PageSize: req.PageSize,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}
	list := orderRes.List
	b := make([]any, len(list))
	for i := range list {
		b[i] = list[i]
	}
	return pages.New(b, req.PageNo, req.PageSize, orderRes.Total), nil
}

func (l *OrderLogic) AddOrder(req *types.ExchangeReq) (string, error) {

	userId := l.ctx.Value("userId").(int64)
	if !req.OrderValid() {
		return "", errors.New("参数传递错误")
	}
	orderResp, err := l.svcCtx.OrderRpc.Add(l.ctx, &order.OrderReq{
		Symbol:    req.Symbol,
		UserId:    userId,
		Direction: req.Direction,
		Type:      req.Type,
		Price:     req.Price,
		Amount:    req.Amount,
	})
	if err != nil {
		logx.Errorw("OrderRpc-AddOrder-ERROR", logx.Field("err", err))
		return "", err
	}
	return orderResp.OrderId, nil

}
