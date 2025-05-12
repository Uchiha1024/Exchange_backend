package domain

import (
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueueDomain struct {
	kafkaClient *database.KafkaClient
}

func NewQueueDomain(kafkaClient *database.KafkaClient) *QueueDomain {
	return &QueueDomain{
		kafkaClient: kafkaClient,
	}
}

const KLINE1MTopic = "kline_1m"
const BtcTransactionTopic = "BTC_TRANSACTION"

func (d *QueueDomain) Send1mKline(data []string, symbol string) {
	klineData := model.NewKline(data, "1m")
	bytes, _ := json.Marshal(klineData)
	msg := database.KafkaData{
		Topic: KLINE1MTopic,
		Data:  bytes,
		Key:   []byte(symbol),
	}
	d.kafkaClient.Send(msg)
	logx.Info("=================发送数据成功==============")
}


func (d *QueueDomain) SendRecharge(value float64, address string, time int64) {
	data := make(map[string]any)
	data["value"] = value
	data["address"] = address
	data["time"] = time
	data["type"] = model.RECHARGE
	data["symbol"] = "BTC"
	marshal, _ := json.Marshal(data)
	msg := database.KafkaData{
		Topic: BtcTransactionTopic,
		Data:  marshal,
		Key:   []byte(address),
	}
	d.kafkaClient.Send(msg)
}