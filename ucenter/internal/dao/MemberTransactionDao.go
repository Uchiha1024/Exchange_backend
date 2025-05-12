package dao

import (
	"context"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
	"mscoin-common/tools"
	"ucenter/internal/model"

	"gorm.io/gorm"
)

type MemberTransactionDao struct {
	conn *gorms.GormConn
}

func NewMemberTransactionDao(db *msdb.MsDB) *MemberTransactionDao {
	return &MemberTransactionDao{
		conn: gorms.New(db.Conn),
	}
}

func (d *MemberTransactionDao) FindTransaction(
	ctx context.Context,
	pageNo int,
	pageSize int,
	memberId int64,
	startTime string,
	endTime string,
	symbol string,
	transactionType string) (list []*model.MemberTransaction, total int64, err error) {
	session := d.conn.Session(ctx)
	db := session.Model(&model.MemberTransaction{}).Where("member_id = ?", memberId)
	if transactionType != "" {
		db = db.Where("type = ?", tools.ToInt64(transactionType))
	}
	if symbol != "" {
		db = db.Where("symbol = ?", symbol)
	}
	if startTime != "" && endTime != "" {
		sTime := tools.ToMill(startTime)
		eTime := tools.ToMill(endTime)
		db.Where("create_time >= ? and create_time <= ?", sTime, eTime)
	}
	offset := (pageNo - 1) * pageSize
	db.Count(&total)
	db.Order("create_time desc").Limit(pageSize).Offset(offset).Find(&list)
	err = db.Error
	return

}

func (d *MemberTransactionDao) Save(ctx context.Context, transaction *model.MemberTransaction) error {
	session := d.conn.Session(ctx)
	err := session.Save(&transaction).Error
	return err
}

func (d *MemberTransactionDao) FindByAmountAndTime(
	ctx context.Context,
	address string,
	value float64,
	time int64) (mt *model.MemberTransaction, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.MemberTransaction{}).
		Where("address=? and amount=? and create_time=?", address, value, time).
		Limit(1).Take(&mt).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
