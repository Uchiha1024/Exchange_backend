package domain

import (
	"context"
	"exchange/internal/dao"
	"exchange/internal/model"
	"exchange/internal/repo"
	"mscoin-common/msdb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderDomain struct {
	orderRepo repo.ExchangeOrderRepo
}

func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{
		orderRepo: dao.NewExchangeOrderDao(db),
	}
}

func (d *ExchangeOrderDomain) FindOrderHistory(ctx context.Context, symbol string, page int64, size int64, memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderHistory(ctx, symbol, page, size, memberId)
	if err != nil {
		logx.Errorw("Domain-FindOrderHistory", logx.Field("error", err))
		return nil, 0, err
	}
	voList := make([]*model.ExchangeOrderVo, 0)
	for i,v := range list {
		voList[i] = v.ToVo()
		
	}
	return voList, total, nil
}

func (d *ExchangeOrderDomain) FindOrderCurrent(ctx context.Context, symbol string, page int64, size int64, memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderCurrent(ctx, symbol, page, size, memberId)
	if err != nil {
		logx.Errorw("Domain-FindOrderCurrent", logx.Field("error", err))
	}
	voList := make([]*model.ExchangeOrderVo, 0)
	for i,v := range list {
		voList[i] = v.ToVo()
		
	}
	return voList, total, nil
}