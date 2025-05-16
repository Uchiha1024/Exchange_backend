package dao

import (
	"context"
	"market/internal/model"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CoinDao struct {
	conn *gorms.GormConn
}

func NewCoinDao(db *msdb.MsDB) *CoinDao {
	return &CoinDao{
		conn: gorms.New(db.Conn),
	}
}

func (d *CoinDao) FindByUnit(ctx context.Context, unit string) (*model.Coin, error) {

	session := d.conn.Session(ctx)
	coin := &model.Coin{}
	err := session.Model(&model.Coin{}).Where("unit = ?", unit).Find(coin).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		logx.Errorf("query findCoinInfo by uint not found error: %v", err)
		return nil, nil
	}
	if err != nil {
		logx.Errorf("query findCoinInfo by uint  error: %v", err)
		return nil, err
	}
	return coin, nil

}


func (d *CoinDao) FindAll(ctx context.Context) (list []*model.Coin, err error) {
	session := d.conn.Session(ctx)
	err = session.Model(&model.Coin{}).Find(&list).Error
	if err != nil {
		logx.Errorf("query findAllCoin error: %v", err)
		return nil, err
	}
	return 
}
