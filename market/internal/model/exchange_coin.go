package model

type ExchangeCoin struct {
	Id               int64   `gorm:"column:id"`
	Symbol           string  `gorm:"column:symbol"`             // 交易币种名称，格式：BTC/USDT
	BaseCoinScale    int64   `gorm:"column:base_coin_scale"`    // 基币小数精度
	BaseSymbol       string  `gorm:"column:base_symbol"`        // 结算币种符号，如USDT
	CoinScale        int64   `gorm:"column:coin_scale"`         // 交易币小数精度
	CoinSymbol       string  `gorm:"column:coin_symbol"`        // 交易币种符号
	Enable           int64   `gorm:"column:enable"`             // 状态，1：启用，2：禁止
	Fee              float64 `gorm:"column:fee"`                // 交易手续费
	Sort             int64   `gorm:"column:sort"`               // 排序，从小到大
	EnableMarketBuy  int64   `gorm:"column:enable_market_buy"`  // 是否启用市价买
	EnableMarketSell int64   `gorm:"column:enable_market_sell"` // 是否启用市价卖
	MinSellPrice     float64 `gorm:"column:min_sell_price"`     // 最低挂单卖价
	Flag             int64   `gorm:"column:flag"`               // 标签位，用于推荐，排序等,默认为0，1表示推荐
	MaxTradingOrder  int64   `gorm:"column:max_trading_order"`  // 最大允许同时交易的订单数，0表示不限制
	MaxTradingTime   int64   `gorm:"column:max_trading_time"`   // 委托超时自动下架时间，单位为秒，0表示不过期
	MinTurnover      float64 `gorm:"column:min_turnover"`       // 最小挂单成交额
	ClearTime        int64   `gorm:"column:clear_time"`         // 清盘时间
	EndTime          int64   `gorm:"column:end_time"`           // 结束时间
	Exchangeable     int64   `gorm:"column:exchangeable"`       //  是否可交易
	MaxBuyPrice      float64 `gorm:"column:max_buy_price"`      // 最高买单价
	MaxVolume        float64 `gorm:"column:max_volume"`         // 最大下单量
	MinVolume        float64 `gorm:"column:min_volume"`         // 最小下单量
	PublishAmount    float64 `gorm:"column:publish_amount"`     //  活动发行数量
	PublishPrice     float64 `gorm:"column:publish_price"`      //  分摊发行价格
	PublishType      int64   `gorm:"column:publish_type"`       // 发行活动类型 1:无活动,2:抢购发行,3:分摊发行
	RobotType        int64   `gorm:"column:robot_type"`         // 机器人类型
	StartTime        int64   `gorm:"column:start_time"`         // 开始时间
	Visible          int64   `gorm:"column:visible"`            //  前台可见状态
	Zone             int64   `gorm:"column:zone"`               // 交易区域
}

func (*ExchangeCoin) TableName() string {
	return "exchange_coin"
}
