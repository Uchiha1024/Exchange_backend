package domain

import (
	"context"
	"errors"
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

func (d *ExchangeCoinDomain) FindBySymbol(ctx context.Context, symbol string) (*model.ExchangeCoin, error) {
	coin, err := d.ExchangeCoinRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		logx.Errorf("query exchange coin by symbol error: %v", err)
		return nil, err
	}
	if coin == nil {
		return nil, errors.New("交易对不存在")
	}
	return coin, nil

}
