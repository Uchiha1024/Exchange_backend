// Code geneorderd by goctl. DO NOT EDIT.
// Source: order.proto

package eclient

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"grpc-common/exchange/types/order"
)

type (
	OrderReq            = order.OrderReq
	OrderRes            = order.OrderRes
	AddOrderRes         = order.AddOrderRes
	ExchangeOrderOrigin = order.ExchangeOrderOrigin
	CancelOrderRes      = order.CancelOrderRes

	Order interface {
		FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
		FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error)
		Add(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*AddOrderRes, error)
		FindByOrderId(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*ExchangeOrderOrigin, error)
		CancelOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*CancelOrderRes, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func (d *defaultOrder) Add(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*AddOrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.Add(ctx, in, opts...)
}

func (d *defaultOrder) FindByOrderId(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*ExchangeOrderOrigin, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.FindByOrderId(ctx, in, opts...)
}

func (d *defaultOrder) CancelOrder(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*CancelOrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.CancelOrder(ctx, in, opts...)
}

func (d *defaultOrder) FindOrderHistory(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.FindOrderHistory(ctx, in, opts...)
}

func (d *defaultOrder) FindOrderCurrent(ctx context.Context, in *OrderReq, opts ...grpc.CallOption) (*OrderRes, error) {
	client := order.NewOrderClient(d.cli.Conn())
	return client.FindOrderCurrent(ctx, in, opts...)
}

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}
