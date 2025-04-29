package domain

import (
	"context"
	"encoding/json"

	"exchange/internal/database"
	"exchange/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaDomain struct {
	cli         *database.KafkaClient
	orderDomain *ExchangeOrderDomain
}

func NewKafkaDomain(cli *database.KafkaClient, orderDomain *ExchangeOrderDomain) *KafkaDomain {
	kafka := &KafkaDomain{
		cli:         cli,
		orderDomain: orderDomain,
	}

	// 启动一个协程，监听kafka的消息

	go kafka.WaitAddOrderResult()
	return kafka
}

// 发送订单消息

func (k *KafkaDomain) SendOrder(topic string, userId int64, orderId string, money float64, symbol string, direction int, baseSymbol string, coinSymbol string) error {
	m := make(map[string]any)
	m["userId"] = userId
	m["orderId"] = orderId
	m["money"] = money
	m["symbol"] = symbol
	m["direction"] = direction
	m["baseSymbol"] = baseSymbol
	m["coinSymbol"] = coinSymbol
	kafaData, _ := json.Marshal(m)
	data := database.KafkaData{
		Topic: topic,
		Key:   []byte(orderId),
		Data:  kafaData,
	}
	// SendSync（同步发送）
	// 同步发送消息到 Kafka。
	// 方法会阻塞，直到消息被 Kafka 确认（即 broker 返回 ack）。
	// 你可以直接拿到发送结果（成功/失败），可以做错误处理和重试。
	// 适合对可靠性要求高、需要确认消息已被写入 Kafka 的场景。
	err := k.cli.SendSync(data)
	logx.Info("创建订单，发消息成功,orderId=" + orderId)
	return err
}

type OrderResult struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

// 监听订单消息 用于修改订单状态
func (k *KafkaDomain) WaitAddOrderResult() {
	cli := k.cli.StartRead("exchange_order_init_complete_trading")
	for {
		kafkaData:= cli.Read()
		logx.Info("读取exchange_order_init_complete_trading 消息成功 orderId:" + string(kafkaData.Key))
		var orderResult OrderResult
		json.Unmarshal(kafkaData.Data, &orderResult)
		// 先根据orderId查询订单
		exchangeOrder,err   := k.orderDomain.FindByOrderId(context.Background(), orderResult.OrderId)
		if err != nil {
			logx.Error(err)
			// 更新订单状态，取消订单
			err := k.orderDomain.UpdateStatusCancel(context.Background(), exchangeOrder.OrderId)
			if err != nil {
				logx.Error("更新订单状态失败，orderId=" + exchangeOrder.OrderId)
				// 重新方法kafka
				k.cli.RPut(kafkaData)
				continue
			}
			if exchangeOrder == nil {
				logx.Error("订单不存在，orderId=" + orderResult.OrderId)
				continue
			}
			if exchangeOrder.Status != model.Init{
				logx.Error("订单已经被处理过")
				continue
			}
			err = k.orderDomain.UpdateOrderStatusTrading(context.Background(), orderResult.OrderId)
			if err != nil {
				logx.Error(err)
				k.cli.RPut(kafkaData)
				continue
			}
		}
	}
}
