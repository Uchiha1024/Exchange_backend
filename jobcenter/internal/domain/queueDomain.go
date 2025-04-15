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
