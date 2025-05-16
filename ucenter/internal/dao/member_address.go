package dao

import (
	"context"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
	"ucenter/internal/model"
)

type MemberAddressDao struct {
	conn *gorms.GormConn
}

func (m *MemberAddressDao) FindByMemIdAndCoinId(ctx context.Context, memId int64, coinId int64) (list []*model.MemberAddress, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberAddress{}).
		Where("member_id=? and coin_id=?", memId, coinId).
		Find(&list).Error
	return
}

func NewMemberAddressDao(db *msdb.MsDB) *MemberAddressDao {
	return &MemberAddressDao{
		conn: gorms.New(db.Conn),
	}
}
