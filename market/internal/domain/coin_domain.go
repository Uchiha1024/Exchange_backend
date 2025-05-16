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

type CoinDomain struct {
	CoinRepo repo.CoinRepo
}

func NewCoinDomain(db *msdb.MsDB) *CoinDomain {
	return &CoinDomain{
		CoinRepo: dao.NewCoinDao(db),
	}
}

func (d *CoinDomain) FindByUnit(ctx context.Context, unit string) (*model.Coin, error) {

	coin, err := d.CoinRepo.FindByUnit(ctx, unit)
	if err != nil {
		logx.Errorf("query FindByUnit error: %v", err)
		return nil, err
	}
	if coin == nil {
		return nil, errors.New("币种不存在")
	}
	return coin, nil

}


func (d *CoinDomain) FindAll(ctx context.Context) ([]*model.Coin, error) {
	return d.CoinRepo.FindAll(ctx)
}
