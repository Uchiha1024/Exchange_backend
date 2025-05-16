package domain

import (
	"grpc-common/market/mclient"
	"mscoin-common/msdb"
)

type WithdrawDomain struct {
	// withdrawRecordRepo repo.WithdrawRecordRepo
	memberWalletDomain *MemberWalletDomain
	marketRpc          mclient.Market
	BitCoinAddress     string
}

func NewWithdrawDomain(
	db *msdb.MsDB,
	marketRpc mclient.Market,
	BitCoinAddress string) *WithdrawDomain {
	return &WithdrawDomain{
		// withdrawRecordRepo: dao.NewWithdrawRecordDao(db),
		memberWalletDomain: NewMemberWalletDomain(db, nil, nil),
		marketRpc:          marketRpc,
		BitCoinAddress:     BitCoinAddress,
	}
}
