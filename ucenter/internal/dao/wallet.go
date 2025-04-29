package dao

import (
	"context"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
	"ucenter/internal/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MemberWalletDao struct {
	conn *gorms.GormConn
}

func NewMemberWalletDao(db *msdb.MsDB) *MemberWalletDao {
	return &MemberWalletDao{
		conn: gorms.New(db.Conn),
	}
}

func (m *MemberWalletDao) FindByIdAndCoinName(ctx context.Context, memId int64, coinName string) (mw *model.MemberWallet, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.MemberWallet{}).Where("member_id = ? and coin_name = ?", memId, coinName).Take(&mw).Error
	if err == gorm.ErrRecordNotFound {
		logx.Errorf("SQL-FindByIdAndCoinName - RECORD NOT FOUND: %v", err)
		return nil, nil
	}
	return
}

func (m *MemberWalletDao) Save(ctx context.Context, mw *model.MemberWallet) error {
	session := m.conn.Session(ctx)
	return session.Save(mw).Error
}

func (m *MemberWalletDao) UpdateFreeze(ctx context.Context, conn msdb.DbConn, memberId int64, symbol string, money float64) error {
	con := conn.(*gorms.GormConn)
	session := con.Tx(ctx)
	sql := "update member_wallet set balance=balance-?, frozen_balance=frozen_balance+? where member_id=? and coin_name=?"
	err := session.Model(&model.MemberWallet{}).Exec(sql, money, money, memberId,symbol).Error
	return err
}
