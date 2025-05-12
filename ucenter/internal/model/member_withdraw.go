package model

import (
	"github.com/jinzhu/copier"
	"grpc-common/market/types/market"
	"mscoin-common/enum"
	"mscoin-common/tools"
)

type WithdrawRecord struct {
	Id                int64   `gorm:"column:id"`
	MemberId          int64   `gorm:"column:member_id"`
	CoinId            int64   `gorm:"column:coin_id"`
	TotalAmount       float64 `gorm:"column:total_amount"`
	Fee               float64 `gorm:"column:fee"`
	ArrivedAmount     float64 `gorm:"column:arrived_amount"`
	Address           string  `gorm:"column:address"`
	Remark            string  `gorm:"column:remark"`
	TransactionNumber string  `gorm:"column:transaction_number"`
	CanAutoWithdraw   int     `gorm:"column:can_auto_withdraw"`
	IsAuto            int     `gorm:"column:isAuto"`
	Status            int     `gorm:"column:status"`
	CreateTime        int64   `gorm:"column:create_time"`
	DealTime          int64   `gorm:"column:deal_time"`
}

func (*WithdrawRecord) TableName() string {
	return "withdraw_record"
}

const (
	Processing = iota
	Waiting
	Fail
	Success
)

var StatusMap = enum.Enum{
	Processing: "PROCESSING",
	Waiting:    "WAITING",
	Fail:       "FAIL",
	Success:    "SUCCESS",
}

type WithdrawRecordVo struct {
	Id                int64        `json:"id" from:"id"`
	MemberId          int64        `json:"memberId" from:"memberId"`
	Coin              *market.Coin `json:"coinId" from:"coinId"`
	TotalAmount       float64      `json:"totalAmount" from:"totalAmount"`
	Fee               float64      `json:"fee" from:"fee"`
	ArrivedAmount     float64      `json:"arrivedAmount" from:"arrivedAmount"`
	Address           string       `json:"address" from:"address"`
	Remark            string       `json:"remark" from:"remark"`
	TransactionNumber string       `json:"transactionNumber" from:"transactionNumber"`
	CanAutoWithdraw   int          `json:"canAutoWithdraw" from:"canAutoWithdraw"`
	IsAuto            int          `json:"isAuto" from:"isAuto"`
	Status            int          `json:"status" from:"status"`
	CreateTime        string       `json:"createTime" from:"createTime"`
	DealTime          string       `json:"dealTime" from:"dealTime"`
}

func (wr *WithdrawRecord) ToVo(coin *market.Coin) *WithdrawRecordVo {
	var vo WithdrawRecordVo
	copier.Copy(&vo, wr)
	vo.Coin = coin
	vo.CreateTime = tools.ToTimeString(wr.CreateTime)
	vo.DealTime = tools.ToTimeString(wr.DealTime)
	return &vo
}
