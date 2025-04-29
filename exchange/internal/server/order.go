package server

import (
	"context"
	"exchange/internal/logic"
	"exchange/internal/svc"
	"grpc-common/exchange/types/order"
)

type OrderServer struct {
	svcCtx *svc.ServiceContext
	order.UnimplementedOrderServer
}

func NewOrderServer(svcCtx *svc.ServiceContext) *OrderServer {
	return &OrderServer{
		svcCtx: svcCtx,
	}
}

func (e *OrderServer) FindOrderHistory(ctx context.Context, req *order.OrderReq) (*order.OrderRes, error) {
	l := logic.NewExchangeOrderLogic(ctx, e.svcCtx)
	return l.FindOrderHistory(req)
}

func (e *OrderServer) FindOrderCurrent(ctx context.Context, req *order.OrderReq) (*order.OrderRes, error) {
	l := logic.NewExchangeOrderLogic(ctx, e.svcCtx)
	return l.FindOrderCurrent(req)
}


func (e *OrderServer) Add(ctx context.Context, req *order.OrderReq) (*order.AddOrderRes, error) {
	l := logic.NewExchangeOrderLogic(ctx, e.svcCtx)
	return l.AddOrder(req)
}


func (e *OrderServer) FindByOrderId(ctx context.Context, req *order.OrderReq) (*order.ExchangeOrderOrigin, error) {
	l := logic.NewExchangeOrderLogic(ctx, e.svcCtx)
	return l.FindByOrderId(req)
}
func (e *OrderServer) CancelOrder(ctx context.Context, req *order.OrderReq) (*order.CancelOrderRes, error) {
	l := logic.NewExchangeOrderLogic(ctx, e.svcCtx)
	return l.CancelOrder(req)
}