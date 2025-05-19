package consumer

import (
	"context"
	"encoding/json"
	"mscoin-common/msdb"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"
	"ucenter/internal/model"
)

func WithdrawConsumer(kafkaCli *database.KafkaClient, db *msdb.MsDB, address string) {
	//获取到提现记录
	//创建BTC网络交易
	//要将交易发送到BTC网络 这时候经过旷工的打包之后 全球可见
	//创建交易的时候 一定要有手续费
	//UTXO unspend 地址的余额  -> 交易的input
	withdrawDomain := domain.NewWithdrawDomain(db, nil, address)
	for {
		kafkaData := kafkaCli.Read()
		var wr model.WithdrawRecord
		json.Unmarshal(kafkaData.Data, &wr)
		ctx := context.Background()
		err := withdrawDomain.Withdraw(ctx, wr)
		if err != nil {
			kafkaCli.Rput(kafkaData)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
