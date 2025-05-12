package model

type MemberAddress struct {
	Id         int64  `gorm:"column:id"`
	MemberId   int64  `gorm:"column:member_id"`
	CoinId     int64  `gorm:"column:coin_id"`
	Address    string `gorm:"column:address"`
	Remark     string `gorm:"column:remark"`
	Status     int    `gorm:"column:status"` // 0正常 1 非法
	CreateTime int64  `gorm:"column:create_time"`
	DeleteTime int64  `gorm:"column:delete_time"`
}

func (MemberAddress) TableName() string {
	return "member_address"
}
