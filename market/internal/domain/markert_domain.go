package domain

import (
	"context"
	"grpc-common/market/types/market"
	"market/internal/dao"
	"market/internal/database"
	"market/internal/model"
	"market/internal/repo"
	"mscoin-common/op"
	"mscoin-common/tools"
	"time"
)

type MarketDomain struct {
	KlineRepo repo.KlineRepo
}

func NewMarketDomain(mongoClient *database.MongoClient) *MarketDomain {
	return &MarketDomain{
		KlineRepo: dao.NewKlineDao(mongoClient.Db),
	}
}

func (d *MarketDomain) SymbolThumbTrend(coins []*model.ExchangeCoin) []*market.CoinThumb {

	coinThumbs := make([]*market.CoinThumb, len(coins))
	for i, v := range coins {
		from := tools.ZeroTime()
		end := time.Now().UnixMilli()
		klines, err := d.KlineRepo.FindBySymbolTime(context.Background(), v.Symbol, "1H", from, end, "")
		if err != nil {
			coinThumbs[i] = model.DefaultCoinThumb(v.Symbol)
			continue
		}
		length := len(klines)
		if length <= 0 {
			coinThumbs[i] = model.DefaultCoinThumb(v.Symbol)
			continue
		}
		//降序排列 0 最新数据 length-1 今天最开始的数据
		//构建趋势
		trend := make([]float64, length)
		var high float64 = 0
		var low float64 = klines[0].LowestPrice
		var volumes float64 = 0
		var turnover float64 = 0
		for i := length - 1; i >= 0; i-- {
			trend[i] = klines[i].ClosePrice
			highestPrice := klines[i].HighestPrice
			if highestPrice > high {
				high = highestPrice
			}
			lowPrice := klines[i].LowestPrice
			if lowPrice < low {
				low = lowPrice
			}
			volumes = op.AddN(volumes, klines[i].Volume, 8)
			turnover = op.AddN(turnover, klines[i].Turnover, 8)
		}
		newKline := klines[0]
		oldKline := klines[length-1]
		thumb := newKline.ToCoinThumb(v.Symbol, oldKline)
		thumb.Trend = trend
		thumb.High = high
		thumb.Low = low
		thumb.Volume = volumes
		thumb.Turnover = turnover
		coinThumbs[i] = thumb

	}

	return coinThumbs

}
