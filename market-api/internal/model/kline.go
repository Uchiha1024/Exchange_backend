package model

import (
	"grpc-common/market/types/market"
	"mscoin-common/op"
)

type Kline struct {
	Period       string  `bson:"period,omitempty" json:"period"`
	OpenPrice    float64 `bson:"openPrice,omitempty" json:"openPrice"`
	HighestPrice float64 `bson:"highestPrice,omitempty" json:"highestPrice"`
	LowestPrice  float64 `bson:"lowestPrice,omitempty" json:"lowestPrice"`
	ClosePrice   float64 `bson:"closePrice,omitempty" json:"closePrice"`
	Time         int64   `bson:"time,omitempty" json:"time"`
	Count        float64 `bson:"count,omitempty" json:"count"`       //成交笔数
	Volume       float64 `bson:"volume,omitempty" json:"volume"`     //成交量
	Turnover     float64 `bson:"turnover,omitempty" json:"turnover"` //成交额
}

func (k *Kline) ToCoinThumb(symbol string, ct *market.CoinThumb) *market.CoinThumb {
	isSame := false
	if ct.Symbol == symbol && ct.DateTime == k.Time {
		//认为这是同一个数据
		isSame = true
	}
	if !isSame {
		ct.Close = k.ClosePrice
		ct.Open = k.OpenPrice
		if ct.High < k.HighestPrice {
			ct.High = k.HighestPrice
		}
		ct.Low = k.LowestPrice
		if ct.Low > k.LowestPrice {
			ct.Low = k.LowestPrice
		}
		ct.Zone = 0
		ct.Volume = op.AddN(k.Volume, ct.Volume, 8)
		ct.Turnover = op.AddN(k.Turnover, ct.Turnover, 8)
		ct.Change = k.LowestPrice - ct.Close
		ct.Chg = op.MulN(op.DivN(ct.Change, ct.Close, 5), 100, 3)
		ct.UsdRate = k.ClosePrice
		ct.BaseUsdRate = 1
		ct.Trend = append(ct.Trend, k.ClosePrice)
		ct.DateTime = k.Time
	}
	return ct
}

func (k *Kline) InitCoinThumb(symbol string) *market.CoinThumb {
	ct := &market.CoinThumb{}
	ct.Symbol = symbol
	ct.Close = k.ClosePrice
	ct.Open = k.OpenPrice
	ct.High = k.HighestPrice
	ct.Volume = k.Volume
	ct.Turnover = k.Turnover
	ct.Low = k.LowestPrice
	ct.Zone = 0
	ct.Change = 0
	ct.Chg = 0.0
	ct.UsdRate = k.ClosePrice
	ct.BaseUsdRate = 1
	ct.Trend = make([]float64, 0)
	ct.DateTime = k.Time
	return ct
}

type CoinThumb struct {
	Symbol       string    `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol"`
	Open         float64   `protobuf:"fixed64,2,opt,name=open,proto3" json:"open"`
	High         float64   `protobuf:"fixed64,3,opt,name=high,proto3" json:"high"`
	Low          float64   `protobuf:"fixed64,4,opt,name=low,proto3" json:"low"`
	Close        float64   `protobuf:"fixed64,5,opt,name=close,proto3" json:"close"`
	Chg          float64   `protobuf:"fixed64,6,opt,name=chg,proto3" json:"chg"`
	Change       float64   `protobuf:"fixed64,7,opt,name=change,proto3" json:"change"`
	Volume       float64   `protobuf:"fixed64,8,opt,name=volume,proto3" json:"volume"`
	Turnover     float64   `protobuf:"fixed64,9,opt,name=turnover,proto3" json:"turnover"`
	LastDayClose float64   `protobuf:"fixed64,10,opt,name=lastDayClose,proto3" json:"lastDayClose"`
	UsdRate      float64   `protobuf:"fixed64,11,opt,name=usdRate,proto3" json:"usdRate"`
	BaseUsdRate  float64   `protobuf:"fixed64,12,opt,name=baseUsdRate,proto3" json:"baseUsdRate"`
	Zone         float64   `protobuf:"fixed64,13,opt,name=zone,proto3" json:"zone"`
	DateTime     int64     `protobuf:"varint,14,opt,name=dateTime,proto3" json:"dateTime"`
	Trend        []float64 `protobuf:"fixed64,15,rep,packed,name=trend,proto3" json:"trend"`
}
