// Package processor 实现了一个完整的数字货币交易撮合引擎
// 包含以下核心功能：
// 1. 支持市价单和限价单的撮合
// 2. 维护买卖盘口信息
// 3. 实现价格优先、时间优先的撮合规则
// 4. 提供并发安全的订单处理
// 5. 支持订单状态管理和通知机制
package processor

import (
	"context"
	"encoding/json"
	"exchange/internal/database"
	"exchange/internal/domain"
	"exchange/internal/model"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"mscoin-common/msdb"
	"mscoin-common/op"
	"sort"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// CoinTradeFactory 交易引擎工厂
// 用于管理和创建不同交易对的撮合引擎实例
// 使用 map 存储不同交易对的撮合引擎，key 为交易对符号（如 BTC/USDT）
type CoinTradeFactory struct {
	tradeMap map[string]*CoinTrade // 存储不同交易对的撮合引擎实例
	mux      sync.RWMutex          // 读写锁，保护 tradeMap 的并发访问
}

// NewCoinTradeFactory 创建新的交易引擎工厂
// 返回一个初始化的 CoinTradeFactory 实例
func NewCoinTradeFactory() *CoinTradeFactory {
	return &CoinTradeFactory{
		tradeMap: make(map[string]*CoinTrade),
	}
}

// Init 初始化交易引擎工厂
// 从市场服务获取所有可见的交易对，并为每个交易对创建撮合引擎
// marketRpc: 市场服务客户端
// client: Kafka客户端，用于发送交易消息
// db: 数据库连接，用于持久化交易数据
func (c *CoinTradeFactory) Init(marketRpc mclient.Market, client *database.KafkaClient, db *msdb.MsDB) {
	ctx := context.Background()
	exchangeCoinRes, err := marketRpc.FindExchangeCoinVisible(ctx, &market.MarketReq{})
	if err != nil {
		logx.Error(err)
		return
	}
	for _, v := range exchangeCoinRes.List {
		c.AddCoinTrade(v.Symbol, NewCoinTrade(v.Symbol, client, db))
	}
}

// NewCoinTrade 创建新的交易对撮合引擎
// symbol: 交易对符号，如 "BTC/USDT"
// cli: Kafka客户端，用于发送交易消息
// db: 数据库连接，用于持久化交易数据
func NewCoinTrade(symbol string, cli *database.KafkaClient, db *msdb.MsDB) *CoinTrade {
	c := &CoinTrade{
		symbol:      symbol,
		kafkaClient: cli,
		db:          db,
	}
	c.init()
	return c
}

// init 初始化交易对撮合引擎
// 创建买卖盘口和限价队列
func (t *CoinTrade) init() {
	t.buyTradePlate = NewTradePlate(t.symbol, model.BUY)
	t.sellTradePlate = NewTradePlate(t.symbol, model.SELL)
	t.buyLimitQueue = &LimitPriceQueue{}
	t.sellLimitQueue = &LimitPriceQueue{}
	t.initData()
}

// CoinTrade 单个交易对的撮合引擎
// 负责处理特定交易对的所有订单撮合逻辑
type CoinTrade struct {
	symbol          string                // 交易对符号，如 "BTC/USDT"
	buyMarketQueue  TradeTimeQueue        // 市价买单队列，按时间排序
	bmMux           sync.RWMutex          // 市价买单队列的读写锁
	sellMarketQueue TradeTimeQueue        // 市价卖单队列，按时间排序
	smMux           sync.RWMutex          // 市价卖单队列的读写锁
	buyLimitQueue   *LimitPriceQueue      // 限价买单队列，按价格从高到低排序
	sellLimitQueue  *LimitPriceQueue      // 限价卖单队列，按价格从低到高排序
	buyTradePlate   *TradePlate           // 买盘盘口信息，显示当前可成交的买单
	sellTradePlate  *TradePlate           // 卖盘盘口信息，显示当前可成交的卖单
	kafkaClient     *database.KafkaClient // Kafka客户端，用于发送交易消息
	db              *msdb.MsDB            // 数据库连接，用于持久化交易数据
}

// TradeTimeQueue 基于时间的订单队列
// 用于市价单的排序，按照订单提交时间升序排列
type TradeTimeQueue []*model.ExchangeOrder

// Len 返回队列长度
func (t TradeTimeQueue) Len() int {
	return len(t)
}

// Less 定义市价单的排序规则：按时间升序
// 确保先提交的市价单优先成交
func (t TradeTimeQueue) Less(i, j int) bool {
	return t[i].Time < t[j].Time
}

// Swap 交换队列中的两个订单
func (t TradeTimeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// LimitPriceQueue 限价单队列
// 用于限价单的排序和管理，包含价格和对应价格的订单列表
type LimitPriceQueue struct {
	mux  sync.RWMutex // 保护队列并发访问的读写锁
	list TradeQueue   // 按价格排序的订单队列
}

// LimitPriceMap 价格档位映射
// 记录特定价格档位的所有订单
type LimitPriceMap struct {
	price float64                // 价格档位
	list  []*model.ExchangeOrder // 该价格档位的所有订单
}

// TradeQueue 价格队列
// 用于限价单的排序，按照价格降序排列（买单）或升序排列（卖单）
type TradeQueue []*LimitPriceMap

// Len 返回队列长度
func (t TradeQueue) Len() int {
	return len(t)
}

// Less 定义限价单的排序规则：按价格降序
// 买单队列：价格高的优先成交
// 卖单队列：价格低的优先成交
func (t TradeQueue) Less(i, j int) bool {
	return t[i].price > t[j].price
}

// Swap 交换队列中的两个价格档位
func (t TradeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// TradePlate 交易盘口
// 用于维护和展示当前市场的买卖盘深度信息
// 包含价格档位、数量、方向等信息
type TradePlate struct {
	Items     []*TradePlateItem `json:"items"` // 盘口档位信息列表
	Symbol    string            // 交易对符号，如 "BTC/USDT"
	direction int               // 方向：1-买盘，2-卖盘
	maxDepth  int               // 最大深度，控制显示多少档
	mux       sync.RWMutex      // 读写锁，保护并发访问
}

// TradePlateItem 盘口档位信息
// 记录每个价格档位的价格和数量
type TradePlateItem struct {
	Price  float64 `json:"price"`  // 价格档位
	Amount float64 `json:"amount"` // 该价格档位的总数量
}

// GetItems 获取盘口信息
// 返回当前盘口的所有价格档位信息
func (p *TradePlate) GetItems() []*TradePlateItem {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.Items
}

// Clear 清空盘口信息
// 用于重置或初始化盘口
func (p *TradePlate) Clear() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.Items = make([]*TradePlateItem, 0)
}

// UpdateAmount 更新指定价格档位的数量
// 用于订单成交后更新盘口数量
// price: 价格档位
// amount: 要更新的数量
func (p *TradePlate) UpdateAmount(price float64, amount float64) {
	p.mux.Lock()
	defer p.mux.Unlock()

	for _, v := range p.Items {
		if v.Price == price {
			v.Amount = op.FloorFloat(v.Amount-amount, 8)
			// 如果数量为0，从盘口中移除该价格档位
			if v.Amount <= 0 {
				// TODO: 实现移除逻辑
			}
			return
		}
	}
}

// NewTradePlate 创建新的交易盘口
// symbol: 交易对符号
// direction: 方向（1-买盘，2-卖盘）
func NewTradePlate(symbol string, direction int) *TradePlate {
	return &TradePlate{
		Symbol:    symbol,
		direction: direction,
		maxDepth:  100,
	}
}

// TradePlateResult 盘口查询结果
// 包含盘口的方向、最大/最小数量、最高/最低价格等信息
type TradePlateResult struct {
	Direction    string            `json:"direction"`    // 方向（买/卖）
	MaxAmount    float64           `json:"maxAmount"`    // 最大数量
	MinAmount    float64           `json:"minAmount"`    // 最小数量
	HighestPrice float64           `json:"highestPrice"` // 最高价格
	LowestPrice  float64           `json:"lowestPrice"`  // 最低价格
	Symbol       string            `json:"symbol"`       // 交易对符号
	Items        []*TradePlateItem `json:"items"`        // 盘口档位信息列表
}

// AllResult 获取完整的盘口信息
// 返回包含所有档位的盘口信息
func (p *TradePlate) AllResult() *TradePlateResult {
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items
	return result
}

// Result 获取指定数量的盘口信息
// num: 要获取的档位数量
func (p *TradePlate) Result(num int) *TradePlateResult {
	if num > len(p.Items) {
		num = len(p.Items)
	}
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items[:num]
	return result
}

// getMaxAmount 获取盘口最大数量
func (p *TradePlate) getMaxAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = 0
	for _, v := range p.Items {
		if v.Amount > amount {
			amount = v.Amount
		}
	}
	return amount
}

// getMinAmount 获取盘口最小数量
func (p *TradePlate) getMinAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = p.Items[0].Amount
	for _, v := range p.Items {
		if v.Amount < amount {
			amount = v.Amount
		}
	}
	return amount
}

// getHighestPrice 获取盘口最高价格
func (p *TradePlate) getHighestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = 0
	for _, v := range p.Items {
		if v.Price > price {
			price = v.Price
		}
	}
	return price
}

// getLowestPrice 获取盘口最低价格
func (p *TradePlate) getLowestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = p.Items[0].Price
	for _, v := range p.Items {
		if v.Price < price {
			price = v.Price
		}
	}
	return price
}

// Remove 从盘口移除订单
// order: 要移除的订单
// amount: 要移除的数量
func (p *TradePlate) Remove(order *model.ExchangeOrder, amount float64) {
	for i, v := range p.Items {
		if v.Price == order.Price {
			v.Amount = op.SubFloor(v.Amount, amount, 8)
			if v.Amount <= 0 {
				p.Items = append(p.Items[:i], p.Items[i+1:]...)
			}
			break
		}
	}
}

// AddCoinTrade 添加交易对撮合引擎
// symbol: 交易对符号
// ct: 交易对撮合引擎实例
func (c *CoinTradeFactory) AddCoinTrade(symbol string, ct *CoinTrade) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.tradeMap[symbol] = ct
}

// GetCoinTrade 获取交易对撮合引擎
// symbol: 交易对符号
// 返回对应的撮合引擎实例
func (c *CoinTradeFactory) GetCoinTrade(symbol string) *CoinTrade {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.tradeMap[symbol]
}

func (t *CoinTrade) initData() {
	orderDomain := domain.NewExchangeOrderDomain(t.db)
	//应该去查询对应symbol的订单 将其赋值到coinTrade里面的各个队列中，同时加入买卖盘
	exchangeOrders, err := orderDomain.FindOrderListBySymbol(context.Background(), t.symbol, model.Trading)
	if err != nil {
		logx.Error(err)
		return
	}
	for _, v := range exchangeOrders {
		if v.Type == model.MarketPrice {
			if v.Direction == model.BUY {
				t.bmMux.Lock()
				t.buyMarketQueue = append(t.buyMarketQueue, v)
				t.bmMux.Unlock()
				continue
			}
			if v.Direction == model.SELL {
				t.smMux.Lock()
				t.sellMarketQueue = append(t.sellMarketQueue, v)
				t.smMux.Unlock()
				continue
			}
			//市价单 不进入买卖盘的
		} else if v.Type == model.LimitPrice {
			if v.Direction == model.BUY {
				t.buyLimitQueue.mux.Lock()
				//deal
				isPut := false
				for _, o := range t.buyLimitQueue.list {
					if o.price == v.Price {
						o.list = append(o.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					lpm := &LimitPriceMap{
						price: v.Price,
						list:  []*model.ExchangeOrder{v},
					}
					t.buyLimitQueue.list = append(t.buyLimitQueue.list, lpm)
				}
				t.buyTradePlate.Add(v)
				t.buyLimitQueue.mux.Unlock()
			} else if v.Direction == model.SELL {
				t.sellLimitQueue.mux.Lock()
				//deal
				isPut := false
				for _, o := range t.sellLimitQueue.list {
					if o.price == v.Price {
						o.list = append(o.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					lpm := &LimitPriceMap{
						price: v.Price,
						list:  []*model.ExchangeOrder{v},
					}
					t.sellLimitQueue.list = append(t.sellLimitQueue.list, lpm)
				}
				t.sellTradePlate.Add(v)
				t.sellLimitQueue.mux.Unlock()
			}
		}
	}
	//排序
	sort.Sort(t.buyMarketQueue)
	sort.Sort(t.sellMarketQueue)
	sort.Sort(t.buyLimitQueue.list)                //从高到低
	sort.Sort(sort.Reverse(t.sellLimitQueue.list)) //从低到高
	if len(exchangeOrders) > 0 {
		t.sendTradPlateMsg(t.buyTradePlate)
		t.sendTradPlateMsg(t.sellTradePlate)
	}
}

// Trade 处理新订单
// 根据订单类型（市价/限价）和方向（买/卖）进行撮合
// exchangeOrder: 要处理的订单
func (t *CoinTrade) Trade(exchangeOrder *model.ExchangeOrder) {
	// 根据订单方向选择对应的队列
	var limitPriceList *LimitPriceQueue
	var marketPriceList TradeTimeQueue
	if exchangeOrder.Direction == model.BUY {
		limitPriceList = t.sellLimitQueue
		marketPriceList = t.sellMarketQueue
	} else {
		limitPriceList = t.buyLimitQueue
		marketPriceList = t.buyMarketQueue
	}

	// 根据订单类型进行撮合
	if exchangeOrder.Type == model.MarketPrice {
		// 市价单与限价单撮合
		t.matchMarketPriceWithLP(limitPriceList, exchangeOrder)
	} else {
		// 限价单先与限价单撮合
		t.matchLimitPriceWithLP(limitPriceList, exchangeOrder)
		// 如果未完全成交，继续与市价单撮合
		if exchangeOrder.Status == model.Trading {
			t.matchLimitPriceWithMP(marketPriceList, exchangeOrder)
		}
		// 如果仍未完全成交，加入限价队列
		if exchangeOrder.Status == model.Trading {
			t.addLimitQueue(exchangeOrder)
			if exchangeOrder.Direction == model.BUY {
				t.sendTradPlateMsg(t.buyTradePlate)
			} else {
				t.sendTradPlateMsg(t.sellTradePlate)
			}
		}
	}
}

// matchLimitPriceWithMP 限价单与市价单撮合
// mpList: 市价单队列
// focusedOrder: 当前要撮合的限价单
func (t *CoinTrade) matchLimitPriceWithMP(mpList TradeTimeQueue, focusedOrder *model.ExchangeOrder) {
	var delOrders []string
	for _, matchOrder := range mpList {
		// 跳过自己的订单，防止自成交
		if matchOrder.MemberId == focusedOrder.MemberId {
			continue
		}
		price := focusedOrder.Price
		// 计算可交易的数量
		matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
		if matchAmount <= 0 {
			continue
		}
		focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
		if matchAmount >= focusedAmount {
			// 完全成交
			turnover := op.MulFloor(price, focusedAmount, 8)
			matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
			matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
			if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
			}
			focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
			focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
			focusedOrder.Status = model.Completed
			break
		} else {
			// 部分成交
			turnover := op.MulFloor(price, matchAmount, 8)
			matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
			matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
			matchOrder.Status = model.Completed
			delOrders = append(delOrders, matchOrder.OrderId)
			focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
			focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
			continue
		}
	}
	// 删除已完成的订单
	for _, orderId := range delOrders {
		for index, matchOrder := range mpList {
			if matchOrder.OrderId == orderId {
				mpList = append(mpList[:index], mpList[index+1:]...)
				break
			}
		}
	}
}

// matchLimitPriceWithLP 限价单与限价单撮合
// lpList: 限价单队列
// focusedOrder: 当前要撮合的限价单
// 只有限价单撮合需要发送完成通知
// 市价单的完成通知可能由其他流程统一处理
// 这样可以避免重复发送通知

func (t *CoinTrade) matchLimitPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	var delOrders []string
	buyNotify := false
	sellNotify := false
	var completeOrders []*model.ExchangeOrder

	// 遍历限价队列
	for _, v := range lpList.list {
		for _, matchOrder := range v.list {
			// 跳过自己的订单
			if matchOrder.MemberId == focusedOrder.MemberId {
				continue
			}
			// 检查价格是否满足成交条件
			if model.BUY == focusedOrder.Direction {
				if focusedOrder.Price < matchOrder.Price {
					break
				}
			}
			if model.SELL == focusedOrder.Direction {
				if focusedOrder.Price > matchOrder.Price {
					break
				}
			}
			// 计算可交易数量
			price := matchOrder.Price
			matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}
			focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if matchAmount >= focusedAmount {
				// 完全成交
				turnover := op.MulFloor(price, focusedAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
					completeOrders = append(completeOrders, matchOrder)
				}
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				focusedOrder.Status = model.Completed
				completeOrders = append(completeOrders, focusedOrder)
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focusedAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focusedAmount)
					sellNotify = true
				}
				break
			} else {
				// 部分成交
				turnover := op.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				matchOrder.Status = model.Completed
				completeOrders = append(completeOrders, matchOrder)
				delOrders = append(delOrders, matchOrder.OrderId)
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
	}
	// 删除已完成的订单
	for _, orderId := range delOrders {
		for _, v := range lpList.list {
			for index, matchOrder := range v.list {
				if orderId == matchOrder.OrderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	// 通知盘口更新
	if buyNotify {
		t.sendTradPlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradPlateMsg(t.sellTradePlate)
	}
	for _, v := range completeOrders {
		t.sendCompleteOrder(v)
	}
}

// matchMarketPriceWithLP 市价单与限价单撮合
// lpList: 限价单队列
// focusedOrder: 当前要撮合的市价单
func (t *CoinTrade) matchMarketPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	// 加写锁，保证并发安全
	lpList.mux.Lock()
	defer lpList.mux.Unlock()

	var delOrders []string
	buyNotify := false
	sellNotify := false

	// 遍历限价队列
	for _, v := range lpList.list {
		for _, matchOrder := range v.list {
			// 跳过自己的订单
			if matchOrder.MemberId == focusedOrder.MemberId {
				continue
			}

			// 获取对方订单价格
			price := matchOrder.Price

			// 计算可交易数量
			matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}

			focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)

			// 市价买单需要根据价格换算数量
			if focusedOrder.Direction == model.BUY {
				focusedAmount = op.DivFloor(op.SubFloor(focusedOrder.Amount, focusedOrder.Turnover, 8), price, 8)
			}

			if matchAmount >= focusedAmount {
				// 完全成交
				turnover := op.MulFloor(price, focusedAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.Completed
					delOrders = append(delOrders, matchOrder.OrderId)
				}
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				focusedOrder.Status = model.Completed
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focusedAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focusedAmount)
					sellNotify = true
				}
				break
			} else {
				// 部分成交
				turnover := op.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				matchOrder.Status = model.Completed
				delOrders = append(delOrders, matchOrder.OrderId)
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
	}

	// 删除已完成的订单
	for _, orderId := range delOrders {
		for _, v := range lpList.list {
			for index, matchOrder := range v.list {
				if orderId == matchOrder.OrderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}

	// 如果未完全成交，重新放入队列
	if focusedOrder.Status == model.Trading {
		t.addMarketQueue(focusedOrder)
	}

	// 通知盘口更新
	if buyNotify {
		t.sendTradPlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradPlateMsg(t.sellTradePlate)
	}
}

// addMarketQueue 添加市价单到队列
// order: 要添加的市价单
func (t *CoinTrade) addMarketQueue(order *model.ExchangeOrder) {
	if order.Type != model.MarketPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyMarketQueue = append(t.buyMarketQueue, order)
		sort.Sort(t.buyMarketQueue)
	} else {
		t.sellMarketQueue = append(t.sellMarketQueue, order)
		sort.Sort(t.sellMarketQueue)
	}
}

// addLimitQueue 添加限价单到队列
// order: 要添加的限价单
func (t *CoinTrade) addLimitQueue(order *model.ExchangeOrder) {
	if order.Type != model.LimitPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyLimitQueue.mux.Lock()
		isPut := false
		for _, o := range t.buyLimitQueue.list {
			if o.price == order.Price {
				o.list = append(o.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			lpm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.buyLimitQueue.list = append(t.buyLimitQueue.list, lpm)
		}
		t.buyTradePlate.Add(order)
		t.buyLimitQueue.mux.Unlock()
	} else if order.Direction == model.SELL {
		t.sellLimitQueue.mux.Lock()
		isPut := false
		for _, o := range t.sellLimitQueue.list {
			if o.price == order.Price {
				o.list = append(o.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			lpm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.sellLimitQueue.list = append(t.sellLimitQueue.list, lpm)
		}
		t.sellTradePlate.Add(order)
		t.sellLimitQueue.mux.Unlock()
	}
}

// sendCompleteOrder 发送订单完成通知
// order: 已完成的订单
func (t *CoinTrade) sendCompleteOrder(order *model.ExchangeOrder) {
	if order.Status != model.Completed {
		return
	}
	marshal, _ := json.Marshal(order)
	kafkaData := database.KafkaData{
		Topic: "exchange_order_complete",
		Key:   []byte(t.symbol),
		Data:  marshal,
	}
	for {
		err := t.kafkaClient.SendSync(kafkaData)
		if err != nil {
			logx.Error(err)
			time.Sleep(250 * time.Millisecond)
			continue
		} else {
			break
		}
	}
}

// Add 添加订单到交易盘口
// 根据订单类型和方向更新盘口信息
// order: 要添加的订单
func (p *TradePlate) Add(order *model.ExchangeOrder) {
	// 检查订单方向是否匹配
	if p.direction != order.Direction {
		return
	}

	// 加锁保护并发访问
	p.mux.Lock()
	defer p.mux.Unlock()

	// 市价单不进入买卖盘
	if order.Type == model.MarketPrice {
		return
	}

	// 获取当前盘口大小
	size := len(p.Items)
	if size > 0 {
		// 遍历现有价格档位
		for _, v := range p.Items {
			// 如果找到相同价格档位，更新数量
			if v.Price == order.Price {
				v.Amount = op.FloorFloat(v.Amount+(order.Amount-order.TradedAmount), 8)
				return
			}
		}
	}

	// 如果是新的价格档位，且未超过最大深度限制
	if size < p.maxDepth {
		tpi := &TradePlateItem{
			Amount: op.FloorFloat(order.Amount-order.TradedAmount, 8),
			Price:  order.Price,
		}
		p.Items = append(p.Items, tpi)
	}
}

// sendTradPlateMsg 发送盘口更新消息
// tradePlate: 要发送的盘口信息
func (p *CoinTrade) sendTradPlateMsg(tradePlate *TradePlate) {
	result := tradePlate.Result(24)
	marshal, _ := json.Marshal(result)
	data := database.KafkaData{
		Topic: "exchange_order_trade_plate",
		Key:   []byte(tradePlate.Symbol),
		Data:  marshal,
	}
	err := p.kafkaClient.SendSync(data)
	if err != nil {
		logx.Error(err)
	} else {
		logx.Info("======exchange_order_trade_plate send 成功....==========")
	}
}
