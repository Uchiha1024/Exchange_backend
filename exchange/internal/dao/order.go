package dao

import (
	"context"
	"exchange/internal/model"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (e *ExchangeOrderDao) FindCurrentTradingCount(ctx context.Context, id int64, symbol string, direction int) (total int64, err error) {
	session := e.conn.Session(ctx)
	err = session.Model(&model.ExchangeOrder{}).
		Where("symbol=? and member_id=? and direction=? and status=?", symbol, id, direction, model.Trading).
		Count(&total).Error
	return
}

func (e *ExchangeOrderDao) Save(ctx context.Context, conn msdb.DbConn, order *model.ExchangeOrder) error {
	e.conn = conn.(*gorms.GormConn)
	tx := e.conn.Tx(ctx)
	err := tx.Save(&order).Error
	return err
}

func (e *ExchangeOrderDao) FindOrderByOrderId(ctx context.Context, orderId string) (order *model.ExchangeOrder, err error) {
	session := e.conn.Session(ctx)

	err = session.Model(&model.ExchangeOrder{}).
		Where("order_id=?", orderId).
		First(&order).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (e *ExchangeOrderDao) UpdateStatusCancel(ctx context.Context, orderId string) error {
	session := e.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).
		Where("order_id=?", orderId).
		Update("status", model.Canceled).Error
	return err
}

func (e *ExchangeOrderDao) UpdateOrderStatusTrading(ctx context.Context, orderId string) error {
	session := e.conn.Session(ctx)
	err := session.Model(&model.ExchangeOrder{}).
		Where("order_id=?", orderId).
		Update("status", model.Trading).Error
	return err
}

