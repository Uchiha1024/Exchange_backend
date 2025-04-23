package dao

import (
	"context"
	"exchange/internal/model"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderDao struct {
	conn *gorms.GormConn
}

func NewExchangeOrderDao(db *msdb.MsDB) *ExchangeOrderDao {
	return &ExchangeOrderDao{
		conn: gorms.New(db.Conn),
	}
}

func (d *ExchangeOrderDao) FindOrderHistory(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := d.conn.Session(ctx)

	err = session.Model(&model.ExchangeOrder{}).
		Where("symbol=? and member_id=?", symbol, memberId).
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&list).Error
	if err != nil {
		logx.Errorw("DAO-FindOrderHistory", logx.Field("error", err))
		return
	}
	err = session.Model(&model.ExchangeOrder{}).
		Where("symbol=? and member_id=?", symbol, memberId).
		Count(&total).Error
	if err != nil {
		logx.Errorw("DAO-FindOrderHistory", logx.Field("error", err))
		return
	}

	return
}

func (e *ExchangeOrderDao) FindOrderCurrent(ctx context.Context, symbol string, page int64, size int64, memberId int64) (list []*model.ExchangeOrder, total int64, err error) {
	session := e.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).
		Where("symbol=? and member_id=? and status=?", symbol, memberId, model.Trading).
		Limit(int(size)).
		Offset(int((page - 1) * size)).Find(&list).Error
	err = session.Model(&model.ExchangeOrder{}).
		Where("symbol=? and member_id=? and status=?", symbol, memberId, model.Trading).
		Count(&total).Error
	if err != nil {
		logx.Errorw("DAO-FindOrderCurrent", logx.Field("error", err))
		return
	}
	return
}
