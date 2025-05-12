package domain

import (
	"context"
	"errors"
	"mscoin-common/msdb"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"
)

type MemberTransactionDomain struct {
	memberTransactionRepo repo.MemberTransactionRepo
	memberWalletDomain    *MemberWalletDomain
}

func NewMemberTransactionDomain(db *msdb.MsDB) *MemberTransactionDomain {
	return &MemberTransactionDomain{
		memberTransactionRepo: dao.NewMemberTransactionDao(db),
		memberWalletDomain:    NewMemberWalletDomain(db, nil, nil),
	}
}

func (d *MemberTransactionDomain) FindTransaction(
	ctx context.Context,
	userId int64,
	pageNo int64,
	pageSize int64,
	symbol string,
	startTime string,
	endTime string,
	t string) ([]*model.MemberTransactionVo, int64, error) {
	memberTransactions, total, err := d.memberTransactionRepo.FindTransaction(
		ctx,
		int(pageNo),
		int(pageSize),
		userId,
		startTime,
		endTime,
		symbol,
		t)
	if err != nil {
		return nil, total, err
	}
	voList := make([]*model.MemberTransactionVo, len(memberTransactions))
	for i, v := range memberTransactions {
		voList[i] = v.ToVo()
	}
	return voList, total, nil
}

func (d *MemberTransactionDomain) SaveRecharge(address string, value float64, time int64, t string, symbol string) error {
	time = time * 1000
	ctx := context.Background()
	memberTransaction, err := d.memberTransactionRepo.FindByAmountAndTime(ctx, address, value, time)
	if err != nil {
		return err
	}
	wallet, err := d.memberWalletDomain.FindByAddress(ctx, address)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("address not exist ")
	}
	if memberTransaction == nil {
		transactionType := model.TypeMap.Code(t)
		memberTransaction = &model.MemberTransaction{}
		memberTransaction.MemberId = wallet.MemberId
		memberTransaction.Address = address
		memberTransaction.Type = transactionType
		memberTransaction.CreateTime = time * 1000
		memberTransaction.Amount = value
		memberTransaction.Symbol = symbol
		err := d.memberTransactionRepo.Save(ctx, memberTransaction)
		if err != nil {
			return err
		}
	}
	return nil
}
