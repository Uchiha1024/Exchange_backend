package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-common/exchange/eclient"
	"grpc-common/exchange/types/order"
	"mscoin-common/enum"
	"mscoin-common/msdb"
	"mscoin-common/msdb/tran"
	"mscoin-common/op"
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

type ExchangeOrder struct {
	Id            int64   `gorm:"column:id" json:"id"`
	OrderId       string  `gorm:"column:order_id" json:"orderId"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
	BaseSymbol    string  `gorm:"column:base_symbol" json:"baseSymbol"`
	CanceledTime  int64   `gorm:"column:canceled_time" json:"canceledTime"`
	CoinSymbol    string  `gorm:"column:coin_symbol" json:"coinSymbol"`
	CompletedTime int64   `gorm:"column:completed_time" json:"completedTime"`
	Direction     int     `gorm:"column:direction" json:"direction"`
	MemberId      int64   `gorm:"column:member_id" json:"memberId"`
	Price         float64 `gorm:"column:price" json:"price"`
	Status        int     `gorm:"column:status" json:"status"`
	Symbol        string  `gorm:"column:symbol" json:"symbol"`
	Time          int64   `gorm:"column:time" json:"time"`
	TradedAmount  float64 `gorm:"column:traded_amount" json:"tradedAmount"`
	Turnover      float64 `gorm:"column:turnover" json:"turnover"`
	Type          int     `gorm:"column:type" json:"type"`
	UseDiscount   string  `gorm:"column:use_discount" json:"useDiscount"`
}

// status
const (
	Trading = iota
	Completed
	Canceled
	OverTimed
	Init
)

var StatusMap = enum.Enum{
	Trading:   "TRADING",
	Completed: "COMPLETED",
	Canceled:  "CANCELED",
	OverTimed: "OVERTIMED",
}

// direction
const (
	BUY = iota
	SELL
)

var DirectionMap = enum.Enum{
	BUY:  "BUY",
	SELL: "SELL",
}

// type
const (
	MarketPrice = iota
	LimitPrice
)

var TypeMap = enum.Enum{
	MarketPrice: "MARKET_PRICE",
	LimitPrice:  "LIMIT_PRICE",
}

func ExchangeOrderComplete(redisCli *redis.Redis, cli *database.KafkaClient, db *msdb.MsDB) {
	//先接收消息
	for {
		kafkaData := cli.Read()
		var order *ExchangeOrder
		json.Unmarshal(kafkaData.Data, &order)
		if order == nil {
			continue
		}
		if order.Status != Completed {
			continue
		}
		logx.Info("收到exchange_order_complete_update_success 消息成功:" + order.OrderId)
		walletDomain := domain.NewMemberWalletDomain(db, nil, nil)
		lock := redis.NewRedisLock(redisCli, fmt.Sprintf("order_complete_update_wallet::%d", order.MemberId))
		acquire, err := lock.Acquire()
		if err != nil {
			logx.Error(err)
			logx.Info("有进程已经拿到锁进行处理了")
			continue
		}
		if acquire {
			// BTC/USDT
			ctx := context.Background()
			if order.Direction == BUY {
				baseWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.BaseSymbol)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				coinWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.CoinSymbol)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				if order.Type == MarketPrice {
					//市价买 amount USDT 冻结的钱  order.turnover扣的钱 还回去的钱 amount-order.turnover
					baseWallet.FrozenBalance = op.SubFloor(baseWallet.FrozenBalance, order.Amount, 8)
					baseWallet.Balance = op.AddFloor(baseWallet.Balance, op.SubFloor(order.Amount, order.Turnover, 8), 8)
					coinWallet.Balance = op.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
				} else {
					//限价买 冻结的钱是 order.price*amount  成交了turnover 还回去的钱 order.price*amount-order.turnover
					floor := op.MulFloor(order.Price, order.Amount, 8)
					baseWallet.FrozenBalance = op.SubFloor(baseWallet.FrozenBalance, floor, 8)
					baseWallet.Balance = op.AddFloor(baseWallet.Balance, op.SubFloor(floor, order.Turnover, 8), 8)
					coinWallet.Balance = op.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
				}
				err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
			} else {
				//卖 不管是市价还是限价 都是卖的 BTC  解冻amount 得到的钱是 order.turnover
				coinWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.CoinSymbol)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				baseWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.BaseSymbol)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}

				coinWallet.FrozenBalance = op.SubFloor(coinWallet.FrozenBalance, order.Amount, 8)
				baseWallet.Balance = op.AddFloor(baseWallet.Balance, order.Turnover, 8)
				err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
				if err != nil {
					logx.Error(err)
					cli.Rput(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
			}
			logx.Info("更新钱包成功:" + order.OrderId)
			lock.Release()
		}

	}
}
