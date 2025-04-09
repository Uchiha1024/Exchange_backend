package dao

import (
	"context"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
	"ucenter/model"

	"gorm.io/gorm"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func NewMemberDao(db *msdb.MsDB) *MemberDao {
	return &MemberDao{
		conn: gorms.New(db.Conn),
	}

}

func (m *MemberDao) UpdateLoginCount(ctx context.Context, id int64, step int) error {

	session := m.conn.Session(ctx)
	err := session.Exec("update member set login_count = login_count+? where id=?", step, id).Error
	return err

}

func (m *MemberDao) Save(ctx context.Context, mem *model.Member) error {
	session := m.conn.Session(ctx)
	return session.Save(mem).Error
}

func (m *MemberDao) FindByPhone(context context.Context, phone string) (*model.Member, error) {
	session := m.conn.Session(context)
	var mem model.Member
	err := session.Model(&model.Member{}).Where("mobile_phone = ?", phone).Limit(1).Take(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mem, nil
}
