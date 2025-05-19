package dao

import (
	"context"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
	"ucenter/internal/model"
)

type WithdrawRecordDao struct {
	conn *gorms.GormConn
}

func NewWithdrawRecordDao(db *msdb.MsDB) *WithdrawRecordDao {
	return &WithdrawRecordDao{
		conn: gorms.New(db.Conn),
	}
}

func (m *WithdrawRecordDao) Save(ctx context.Context, record *model.WithdrawRecord) error {
	session := m.conn.Session(ctx)
	err := session.Save(record).Error
	return err
}


func (m *WithdrawRecordDao) UpdateSuccess(ctx context.Context, record model.WithdrawRecord) error {
	session := m.conn.Session(ctx)
	err := session.Model(&model.WithdrawRecord{}).
		Where("id=?", record.Id).
		Updates(map[string]any{"transaction_number": record.TransactionNumber, "status": record.Status, "deal_time": record.DealTime}).
		Error
	return err
}


func (m *WithdrawRecordDao) FindByUserId(ctx context.Context, userId int64, page int64, pageSize int64) (list []*model.WithdrawRecord, total int64, err error) {
	session := m.conn.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&model.WithdrawRecord{}).
		Where("member_id=?", userId).
		Limit(int(pageSize)).
		Offset(int(offset)).Find(&list).Error
	err = session.Model(&model.WithdrawRecord{}).
		Where("member_id=?", userId).
		Count(&total).Error
	return
}