package domain

import (
	"context"
	"market/internal/dao"
	"market/internal/model"
	"market/internal/repo"
	"mscoin-common/msdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeCoinDomain struct {
	ExchangeCoinRepo repo.ExchangeCoinRepo
}

func NewExchangeCoinDomain(db *msdb.MsDB) *ExchangeCoinDomain {
	return &ExchangeCoinDomain{
		ExchangeCoinRepo: dao.NewExchangeCoinDao(db),
	}
}

func (d *ExchangeCoinDomain) FindVisible(ctx context.Context) []*model.ExchangeCoin {
	list, err := d.ExchangeCoinRepo.FindVisible(ctx)
	if err != nil {
		logx.Errorf("query visible exchange coin error: %v", err)
		return []*model.ExchangeCoin{}
	}
	return list

}
