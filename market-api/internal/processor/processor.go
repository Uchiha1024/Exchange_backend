package processor

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"market-api/internal/database"
	"market-api/internal/model"
)

const KLINE1M = "kline_1m"
const KLINE = "kline"
const TRADE = "trade"
const TradePlateTopic = "exchange_order_trade_plate"
const TradePlate = "tradePlate"


// 主题接口（Subject）
type Processor interface {
	GetThumb() any
	Process(data ProcessData)
	AddHandler(h MarketHandler)
}


// 观察者接口（Observer）
type MarketHandler interface {
	HandleTrade(symbol string, data []byte)
	HandleKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb)
	HandleTradePlate(symbol string, tp *model.TradePlateResult)
}


// 具体主题（Concrete Subject）
type DefaultProcessor struct {
	kafkaCli *database.KafkaClient
	handlers []MarketHandler
	thumbMap map[string]*market.CoinThumb
}

func NewDefaultProcessor(kafkaCli *database.KafkaClient) *DefaultProcessor {
	return &DefaultProcessor{
		kafkaCli: kafkaCli,
		handlers: make([]MarketHandler, 0),
		thumbMap: make(map[string]*market.CoinThumb),
	}
}


// 添加观察者
func (d *DefaultProcessor) AddHandler(h MarketHandler) {
	//发送到websocket的服务
	d.handlers = append(d.handlers, h)
}

func (p *DefaultProcessor) Init(marketRpc mclient.Market) {
	p.startReadFromKafka(KLINE1M, KLINE)
	p.startReadTradePlate(TradePlateTopic)
	p.initThumbMap(marketRpc)
}

type ProcessData struct {
	Type string //trade 交易 kline k线
	Key  []byte
	Data []byte
}



func (p *DefaultProcessor) startReadFromKafka(topic string, tp string) {
	//一定要先start 后read
	p.kafkaCli.StartRead(topic)
	go p.dealQueueData(p.kafkaCli, tp)
}

func (p *DefaultProcessor) startReadTradePlate(topic string) {
	cli := p.kafkaCli.StartReadNew(topic)
	go p.dealQueueData(cli, TradePlate)
}


func (p *DefaultProcessor) dealQueueData(cli *database.KafkaClient, tp string) {
	//这就是队列的数据
	for {
		msg := cli.Read()
		data := ProcessData{
			Type: tp,
			Key:  msg.Key,
			Data: msg.Data,
		}
		p.Process(data)
	}

}

// 当数据发生变化时，通知所有观察者
func (d *DefaultProcessor) Process(data ProcessData) {
	if data.Type == KLINE {
		symbol := string(data.Key)
		kline := &model.Kline{}
		json.Unmarshal(data.Data, kline)
		 // 通知所有观察者
		for _, v := range d.handlers {
			v.HandleKLine(symbol, kline, d.thumbMap)
		}
	} else if data.Type == TradePlate {
		symbol := string(data.Key)
		tp := &model.TradePlateResult{}
		json.Unmarshal(data.Data, tp)
		for _, v := range d.handlers {
			v.HandleTradePlate(symbol, tp)
		}
	}
}


func (d *DefaultProcessor) GetThumb() any {
	cs := make([]*market.CoinThumb, len(d.thumbMap))
	i := 0
	for _, v := range d.thumbMap {
		cs[i] = v
		i++
	}
	return cs
}


func (d *DefaultProcessor) initThumbMap(marketRpc mclient.Market) {
	symbolThumbRes, err := marketRpc.FindSymbolThumbTrend(context.Background(),
		&market.MarketReq{})
	if err != nil {
		logx.Info(err)
	} else {
		coinThumbs := symbolThumbRes.List
		for _, v := range coinThumbs {
			d.thumbMap[v.Symbol] = v
		}
	}
}


