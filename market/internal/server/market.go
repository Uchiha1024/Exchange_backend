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

func (s *MarketServer) HistoryKline(ctx context.Context, in *market.MarketReq) (*market.HistoryRes, error) {
	logic := logic.NewMarketLogic(ctx, s.svcCtx)
	return logic.HistoryKline(in)
}


func (e *MarketServer) FindExchangeCoinVisible(ctx context.Context, req *market.MarketReq) (*market.ExchangeCoinRes, error) {
	l := logic.NewMarketLogic(ctx, e.svcCtx)
	return l.FindExchangeCoinVisible(req)
}

func (e *MarketServer) FindAllCoin(ctx context.Context, req *market.MarketReq) (*market.CoinList, error) {
	l := logic.NewMarketLogic(ctx, e.svcCtx)
	return l.FindAllCoin(req)
}

func (e *MarketServer) FindCoinById(ctx context.Context, req *market.MarketReq) (*market.Coin, error) {
	l := logic.NewMarketLogic(ctx, e.svcCtx)
	return l.FindById(req)
}