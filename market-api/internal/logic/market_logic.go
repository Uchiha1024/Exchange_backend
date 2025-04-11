package logic

import (
	"context"
	"grpc-common/market/types/market"
	"market-api/internal/svc"
	"market-api/internal/types"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketLogic {
	return &MarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketLogic) SymbolThumbTrend(req *types.MarketReq) ([]*types.CoinThumbResp, error) {

	

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	symbolThumbResp, err := l.svcCtx.MarketRpc.FindSymbolThumbTrend(ctx, &market.MarketReq{
		Ip: req.Ip,

	})
	if err != nil {
		return nil, err
	}


	var coinThumbResp [] *types.CoinThumbResp
	if err := copier.Copy(&coinThumbResp, symbolThumbResp.List); err != nil {
		return nil, err
	}

	return coinThumbResp, nil


}


