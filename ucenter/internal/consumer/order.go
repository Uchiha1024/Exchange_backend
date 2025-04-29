package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-common/exchange/eclient"
	"grpc-common/exchange/types/order"
	"mscoin-common/msdb"
	"mscoin-common/msdb/tran"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type OrderAdd struct {
	UserId     int64   `json:"userId"`
	OrderId    string  `json:"orderId"`
	Money      float64 `json:"money"`
	Symbol     string  `json:"symbol"`
	Direction  int     `json:"direction"`
	BaseSymbol string  `json:"baseSymbol"`
	CoinSymbol string  `json:"coinSymbol"`
}

func ExchangeOrderAdd(redisCli *redis.Redis, kafkaCli *database.KafkaClient, orderRpc eclient.Order, db *msdb.MsDB) {

	for {
		// 从kafka中读取消息
		kafaData := kafkaCli.Read()
		if kafaData.Data == nil {
			logx.Info("kafka 没有数据")
			continue
		}
		orderId := string(kafaData.Key)
		logx.Infof("kafka 接收到创建订单消息: %s", orderId)
		var orderAdd OrderAdd
		json.Unmarshal(kafaData.Data, &orderAdd)
		if orderAdd.OrderId != orderId {
			logx.Error("订单id不匹配")
			continue
		}
		ctx := context.Background()
		// 查询订单信息
		exchangeOrder, err := orderRpc.FindByOrderId(ctx, &order.OrderReq{
			OrderId: orderId,
		})
		if err != nil {
			logx.Error("查询订单信息失败", err)
			// 取消订单
			cancelOrder(ctx, kafaData, orderId, orderRpc, kafkaCli)
			continue
		}
		if exchangeOrder == nil {
			logx.Error("orderId :" + orderId + " 不存在")
			continue
		}
		//4 init状态
		if exchangeOrder.Status != 4 {
			logx.Error("orderId :" + orderId + " 已经被操作过了")
			continue
		}
		// 添加 Redis 分布式锁
		lock := redis.NewRedisLock(redisCli, "exchange_order::"+fmt.Sprintf("%d::%s", orderAdd.UserId, orderId))
		acquired, err := lock.Acquire()
		if err != nil {
			logx.Error("获取锁失败", err)
			continue
		}
		if acquired {
			//添加事务
			transaction := tran.NewTransaction(db.Conn)
			walletDomain := domain.NewMemberWalletDomain(db, nil, nil)
			err := transaction.Action(func(conn msdb.DbConn) error {
				if orderAdd.Direction == 0 {
					// 买入
					err := walletDomain.Freeze(ctx, conn, orderAdd.UserId, orderAdd.Money, orderAdd.BaseSymbol)
					return err
				} else {
					err := walletDomain.Freeze(ctx, conn, orderAdd.UserId, orderAdd.Money, orderAdd.CoinSymbol)
					return err
				}

			})
			if err != nil {
				cancelOrder(ctx, kafaData, orderId, orderRpc, kafkaCli)
				continue
			}

			//需要将状态 改为trading
			//都完成后 通知订单进行状态变更 需要保证一定发送成功
			for {
				m := make(map[string]any)
				m["userId"] = orderAdd.UserId
				m["orderId"] = orderId
				marshal, _ := json.Marshal(m)
				data := database.KafkaData{
					Topic: "exchange_order_init_complete_trading",
					Key:   []byte(orderId),
					Data:  marshal,
				}
				err := kafkaCli.SendSync(data)
				if err != nil {
					logx.Error(err)
					time.Sleep(250 * time.Millisecond)
					continue
				}
				logx.Info("发送exchange_order_init_complete_trading 消息成功:" + orderId)
				break
			}
			lock.Release()

		}

	}

}

func cancelOrder(ctx context.Context, data database.KafkaData, orderId string, orderRpc eclient.Order, cli *database.KafkaClient) {
	_, err := orderRpc.CancelOrder(ctx, &order.OrderReq{
		OrderId: orderId,
	})
	if err != nil {
		cli.Rput(data)
	}
}
