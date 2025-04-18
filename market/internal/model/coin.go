package model

type Coin struct {
	Id                int     `gorm:"column:id"`
	Name              string  `gorm:"column:name"`
	CanAutoWithdraw   int     `gorm:"column:can_auto_withdraw"`
	CanRecharge       int     `gorm:"column:can_recharge"`
	CanTransfer       int     `gorm:"column:can_transfer"`
	CanWithdraw       int     `gorm:"column:can_withdraw"`
	CnyRate           float64 `gorm:"column:cny_rate"`
	EnableRpc         int     `gorm:"column:enable_rpc"`
	IsPlatformCoin    int     `gorm:"column:is_platform_coin"`
	MaxTxFee          float64 `gorm:"column:max_tx_fee"`
	MaxWithdrawAmount float64 `gorm:"column:max_withdraw_amount"`
	MinTxFee          float64 `gorm:"column:min_tx_fee"`
	MinWithdrawAmount float64 `gorm:"column:min_withdraw_amount"`
	NameCn            string  `gorm:"column:name_cn"`
	Sort              int     `gorm:"column:sort"`
	Status            int     `gorm:"column:status"`
	Unit              string  `gorm:"column:unit"`
	UsdRate           float64 `gorm:"column:usd_rate"`
	WithdrawThreshold float64 `gorm:"column:withdraw_threshold"`
	HasLegal          int     `gorm:"column:has_legal"`
	ColdWalletAddress string  `gorm:"column:cold_wallet_address"`
	MinerFee          float64 `gorm:"column:miner_fee"`
	WithdrawScale     int     `gorm:"column:withdraw_scale"`
	AccountType       int     `gorm:"column:account_type"`
	DepositAddress    string  `gorm:"column:deposit_address"`
	Infolink          string  `gorm:"column:infolink"`
	Information       string  `gorm:"column:information"`
	MinRechargeAmount float64 `gorm:"column:min_recharge_amount"`
}

func (*Coin) TableName() string {
	return "coin"
}
