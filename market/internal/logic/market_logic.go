package logic

import (
	"context"
	"grpc-common/market/types/market"
	"market/internal/domain"
	"market/internal/svc"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarketLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	exchangeCoinDomain *domain.ExchangeCoinDomain
	marketDomain       *domain.MarketDomain
	coinDomain         *domain.CoinDomain
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		svcCtx:             svcCtx,
		exchangeCoinDomain: domain.NewExchangeCoinDomain(svcCtx.Db),
		marketDomain:       domain.NewMarketDomain(svcCtx.MongoClient),
		coinDomain:         domain.NewCoinDomain(svcCtx.Db),
	}
}

func (l *MarketLogic) SymbolThumbTrend(req *market.MarketReq) (*market.SymbolThumbRes, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	coins := l.exchangeCoinDomain.FindVisible(context.Background())
	for i, coin := range coins {
		logx.Infof("coin[%d]: %+v", i, coin)
	}
	coinThumbs := l.marketDomain.SymbolThumbTrend(coins)
	// coinThumbs := make([]*market.CoinThumb, len(coins))
	// for i, v := range coins {
	// 	ct := &market.CoinThumb{}
	// 	ct.Symbol = v.Symbol
	// 	trend := make([]float64, 0)
	// 	for p := 0; p <= 24; p++ {
	// 		trend = append(trend, rand.Float64())
	// 	}
	// 	ct.Trend = trend
	// 	coinThumbs[i] = ct
	// }

	return &market.SymbolThumbRes{
		List: coinThumbs,
	}, nil

}

func (l *MarketLogic) SymbolInfo(in *market.MarketReq) (*market.ExchangeCoin, error) {

	coin, err := l.exchangeCoinDomain.FindBySymbol(l.ctx, in.Symbol)
	if err != nil {
		return nil, err
	}
	exchangeCoin := &market.ExchangeCoin{}
	if err := copier.Copy(exchangeCoin, coin); err != nil {
		logx.Errorf("copy exchange coin to response error: %v", err)
		return nil, err
	}
	return exchangeCoin, nil

}



func (l *MarketLogic) CoinInfo(in *market.MarketReq) (*market.Coin, error) {
	coin, err := l.coinDomain.FindByUnit(l.ctx, in.Unit)
	if err != nil {
		return nil, err
	}
	coinInfo := &market.Coin{}
	if err := copier.Copy(coinInfo, coin); err != nil {
		logx.Errorf("copy coinInfo to response error: %v", err)
		return nil, err
	}
	return coinInfo, nil

}


func (l *MarketLogic) HistoryKline(req *market.MarketReq) (*market.HistoryRes, error) {
	period := "1H"
	if req.Resolution == "60" {
		period = "1H"
	} else if req.Resolution == "30" {
		period = "30m"
	} else if req.Resolution == "15" {
		period = "15m"
	} else if req.Resolution == "5" {
		period = "5m"
	} else if req.Resolution == "1" {
		period = "1m"
	} else if req.Resolution == "1D" {
		period = "1D"
	} else if req.Resolution == "1W" {
		period = "1W"
	} else if req.Resolution == "1M" {
		period = "1M"
	}

	histories,err := l.marketDomain.HistoryKline(l.ctx,req.Symbol,period,req.From,req.To)
	if err != nil {
		return nil,err
	}

	return &market.HistoryRes{
		List: histories,
	}, nil

}


func (l *MarketLogic) FindExchangeCoinVisible(req *market.MarketReq) (*market.ExchangeCoinRes, error) {
	exchangeCoinRes := l.exchangeCoinDomain.FindVisible(l.ctx)
	var list []*market.ExchangeCoin
	copier.Copy(&list, exchangeCoinRes)
	return &market.ExchangeCoinRes{
		List: list,
	}, nil
}
