package server

import (
	"context"
	"grpc-common/market/types/market"

	"market/internal/logic"
	"market/internal/svc"
)

type MarketServer struct {
	svcCtx *svc.ServiceContext
	market.UnimplementedMarketServer
}

func NewMarketServer(svcCtx *svc.ServiceContext) *MarketServer {
	return &MarketServer{
		svcCtx: svcCtx,
		
	}
}

func (s *MarketServer) FindSymbolThumbTrend(ctx context.Context, in *market.MarketReq) (*market.SymbolThumbRes, error) {
	logic := logic.NewMarketLogic(ctx, s.svcCtx)
	return logic.SymbolThumbTrend(in)
}

func (s *MarketServer) FindSymbolInfo(ctx context.Context, in *market.MarketReq) (*market.ExchangeCoin, error) {
	logic := logic.NewMarketLogic(ctx, s.svcCtx)
	return logic.SymbolInfo(in)
}

func (s *MarketServer) FindCoinInfo(ctx context.Context, in *market.MarketReq) (*market.Coin, error) {
	logic := logic.NewMarketLogic(ctx, s.svcCtx)
	return logic.CoinInfo(in)
}
