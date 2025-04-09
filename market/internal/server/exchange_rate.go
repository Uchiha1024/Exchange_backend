package server

import (
	"context"
	"grpc-common/market/types/rate"
	"market/internal/logic"
	"market/internal/svc"
)

type ExchangeRateServer struct {
	svcCtx *svc.ServiceContext
	rate.UnimplementedExchangeRateServer
}

func NewExchangeRateServer(svcCtx *svc.ServiceContext) *ExchangeRateServer {
	return &ExchangeRateServer{
		svcCtx: svcCtx,
		
	}
}

func (s *ExchangeRateServer) UsdRate(ctx context.Context, in *rate.RateReq) (*rate.RateRes, error) {
	logic := logic.NewExchangeRateLogic(ctx, s.svcCtx)
	return logic.UsdRate(in)
}
