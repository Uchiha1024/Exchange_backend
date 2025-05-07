package consumer

import (
	"context"
	"encoding/json"
	"exchange/internal/database"
	"exchange/internal/domain"
	"time"

	// "exchange/internal/domain"
	"exchange/internal/model"
	"exchange/internal/processor"
	"mscoin-common/msdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaConsumer struct {
	cli     *database.KafkaClient
	factory *processor.CoinTradeFactory
	db      *msdb.MsDB
}

func NewKafkaConsumer(cli *database.KafkaClient, factory *processor.CoinTradeFactory, db *msdb.MsDB) *KafkaConsumer {
	return &KafkaConsumer{
		cli:     cli,
		factory: factory,
		db:      db,
	}
}

//消费订单的消息 拿到新创建的订单
//1. 先实现买卖盘的逻辑 买 卖 一旦匹配完成 成交了 成交的价格和数量  就会做为别人的参考 买卖盘也是实时

func (k *KafkaConsumer) Run() {
	orderDomain := domain.NewExchangeOrderDomain(k.db)
	k.orderTrading()
	k.orderComplete(orderDomain)

}

func (k *KafkaConsumer) orderTrading() {
	cli := k.cli.StartRead("exchange_order_trading")
	go k.readOrderTrading(cli)
}



func (k *KafkaConsumer) readOrderTrading(cli *database.KafkaClient) {
	// 解析订单，发送交易引擎
	for {
		kafkaData := cli.Read()
		logx.Info("===== Topic === exchange_order_trading == kafkaData========",string(kafkaData.Data))
		var orderInfo *model.ExchangeOrder
		json.Unmarshal(kafkaData.Data, &orderInfo)
		// 交易对 发送交易引擎
		coinTrade := k.factory.GetCoinTrade(orderInfo.Symbol)
		coinTrade.Trade(orderInfo)

	}
}

func (k *KafkaConsumer) orderComplete(orderDomain *domain.ExchangeOrderDomain) {
	cli := k.cli.StartRead("exchange_order_complete")
	go k.readOrderComplete(cli, orderDomain)
}

func (k *KafkaConsumer) readOrderComplete(cli *database.KafkaClient, orderDomain *domain.ExchangeOrderDomain) {
	for {
		kafkaData := cli.Read()
		logx.Info("===== Topic === exchange_order_complete == kafkaData========",string(kafkaData.Data))
		var orderInfo *model.ExchangeOrder
		json.Unmarshal(kafkaData.Data, &orderInfo)
		// 更新订单状态
		err := orderDomain.UpdateOrderComplete(context.Background(),orderInfo)
		if err != nil {
			logx.Error("===== Topic === exchange_order_complete == kafkaData========",err)
			cli.RPut(kafkaData)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		// 通知钱包更新
		for {
			kafkaData.Topic = "exchange_order_complete_update_success"
			err2 := cli.SendSync(kafkaData)
			if err2 != nil {
				logx.Error("===== Topic === exchange_order_complete_update_success == kafkaData========",err2)
				time.Sleep(200 * time.Millisecond)
				continue
			}
			break
		}
	}
}




