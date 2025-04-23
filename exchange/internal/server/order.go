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
