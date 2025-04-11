package dao

import (
	"context"
	"errors"
	"market/internal/model"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeCoinDao struct {
	conn *gorms.GormConn
}

func NewExchangeCoinDao(db *msdb.MsDB) *ExchangeCoinDao {
	return &ExchangeCoinDao{
		conn: gorms.New(db.Conn),
	}
}

func (d *ExchangeCoinDao) FindVisible(ctx context.Context) ([]*model.ExchangeCoin, error) {
	session := d.conn.Session(ctx)
	var list []*model.ExchangeCoin
	err := session.Model(&model.ExchangeCoin{}).Where("visible = ?", 1).Find(&list).Error
	if err != nil {
		return nil, errors.New("query visible exchange coin error")
	}
	// 添加日志
	// 打印 list
	logx.Info(list)
	logx.Infof("Query result count: %d", len(list))
	return list, nil

}
