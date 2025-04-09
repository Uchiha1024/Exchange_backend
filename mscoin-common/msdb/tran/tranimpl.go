package tran

import (
	"gorm.io/gorm"
	"mscoin-common/msdb"
	"mscoin-common/msdb/gorms"
)

type TransactionImpl struct {
	conn msdb.DbConn
}

func (t *TransactionImpl) Action(f func(conn msdb.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction(db *gorm.DB) *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.New(db),
	}
}
