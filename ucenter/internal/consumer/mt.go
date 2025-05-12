package consumer

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"mscoin-common/msdb"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"
)

type BitCoinTransactionResult struct {
	Value   float64 `json:"value"`
	Time    int64   `json:"time"`
	Address string  `json:"address"`
	Type    string  `json:"type"`
	Symbol  string  `json:"symbol"`
}

func BitCoinTransaction(redisCli *redis.Redis, kafkaCli *database.KafkaClient, db *msdb.MsDB) {
	for {
		kafkaData := kafkaCli.Read()
		var bt BitCoinTransactionResult
		json.Unmarshal(kafkaData.Data, &bt)
		//解析出来数据 调用domain存储到数据库即可
		transactionDomain := domain.NewMemberTransactionDomain(db)
		err := transactionDomain.SaveRecharge(bt.Address, bt.Value, bt.Time, bt.Type, bt.Symbol)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			kafkaCli.Rput(kafkaData)
		}
	}
}
