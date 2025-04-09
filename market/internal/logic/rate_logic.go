package logic

import (
	"context"
	"grpc-common/market/types/rate"
	"market/internal/domain"
	"market/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeRateLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	exchangeDomain *domain.ExchangeRateDomain
}

func NewExchangeRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeRateLogic {
	return &ExchangeRateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeRateLogic) UsdRate(req *rate.RateReq) (*rate.RateRes, error) {

	usdRate := l.exchangeDomain.UsdRate(req.Unit)

	return &rate.RateRes{
		Rate: usdRate,
	}, nil

}
