package model

import (
	"github.com/jinzhu/copier"
	"mscoin-common/enum"
	"mscoin-common/tools"
)

type MemberTransaction struct {
	Id          int64   `gorm:"column:id"`
	Address     string  `gorm:"column:address"`
	Amount      float64 `gorm:"column:amount"`
	CreateTime  int64   `gorm:"column:create_time"`
	Fee         float64 `gorm:"column:fee"`
	Flag        int     `gorm:"column:flag"`
	MemberId    int64   `gorm:"column:member_id"`
	Symbol      string  `gorm:"column:symbol"`
	Type        int     `gorm:"column:type"`
	DiscountFee string  `gorm:"column:discount_fee"`
	RealFee     string  `gorm:"column:real_fee"`
}

func (*MemberTransaction) TableName() string {
	return "member_transaction"
}

const (
	RECHARGE          = iota // 充值
	WITHDRAW                 // 提现
	TRANSFER_ACCOUNTS        //转账
	EXCHANGE                 //币币交易

)

var TypeMap = enum.Enum{
	RECHARGE:          "RECHARGE",
	WITHDRAW:          "WITHDRAW",
	TRANSFER_ACCOUNTS: "TRANSFER_ACCOUNTS",
	EXCHANGE:          "EXCHANGE",
}

type MemberTransactionVo struct {
	Id          int64   `json:"id" from:"id"`
	Address     string  `json:"address" from:"address"`
	Amount      float64 `json:"amount" from:"amount"`
	CreateTime  string  `json:"createTime" from:"createTime"`
	Fee         float64 `json:"fee" from:"fee"`
	Flag        int     `json:"flag" from:"flag"`
	MemberId    int64   `json:"memberId" from:"memberId"`
	Symbol      string  `json:"symbol" from:"symbol"`
	Type        string  `json:"type" from:"type"`
	DiscountFee string  `json:"discountFee" from:"discountFee"`
	RealFee     string  `json:"realFee" from:"realFee"`
}

func (mt *MemberTransaction) ToVo() *MemberTransactionVo {
	vo := &MemberTransactionVo{}
	copier.Copy(vo, mt)
	vo.CreateTime = tools.ToTimeString(mt.CreateTime)
	vo.Type = TypeMap.Value(mt.Type)
	return vo
}
