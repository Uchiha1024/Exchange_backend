package logic

import (
	"context"
	"grpc-common/market/types/market"
	"market/internal/domain"
	"market/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketLogic struct {
	logx.Logger
	ctx                context.Context
	svcCtx             *svc.ServiceContext
	exchangeCoinDomain *domain.ExchangeCoinDomain
	marketDomain       *domain.MarketDomain
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		svcCtx:             svcCtx,
		exchangeCoinDomain: domain.NewExchangeCoinDomain(svcCtx.Db),
		marketDomain:       domain.NewMarketDomain(svcCtx.MongoClient),
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
