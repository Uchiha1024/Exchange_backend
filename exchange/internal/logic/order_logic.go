package logic

import (
	"context"
	"exchange/internal/domain"
	"exchange/internal/svc"
	"grpc-common/exchange/types/order"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeOrderDomain *domain.ExchangeOrderDomain
	// transaction         tran.Transaction
	// kafkaDomain         *domain.KafkaDomain
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	orderDomain := domain.NewExchangeOrderDomain(svcCtx.Db)
	return &ExchangeOrderLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		exchangeOrderDomain: orderDomain,
		// transaction:         tran.NewTransaction(svcCtx.Db.Conn),
		// kafkaDomain:         domain.NewKafkaDomain(svcCtx.KafkaClient, orderDomain),
	}
}

func (l *ExchangeOrderLogic) FindOrderHistory(req *order.OrderReq) (*order.OrderRes, error) {
	voList, total, err := l.exchangeOrderDomain.FindOrderHistory(l.ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory", logx.Field("error", err))
		return nil, err
	}
	var list []*order.ExchangeOrder
	err = copier.Copy(&list, &voList)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory Copier Error", logx.Field("error", err))
		return nil, err
	}
	return &order.OrderRes{
		List: list,
		Total: total,
	}, nil
}


func (l *ExchangeOrderLogic) FindOrderCurrent(req *order.OrderReq) (*order.OrderRes, error) {
	voList, total, err := l.exchangeOrderDomain.FindOrderHistory(l.ctx, req.Symbol, req.Page, req.PageSize, req.UserId)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory", logx.Field("error", err))
		return nil, err
	}
	var list []*order.ExchangeOrder
	err = copier.Copy(&list, &voList)
	if err != nil {
		logx.Errorw("Logic-FindOrderHistory Copier Error", logx.Field("error", err))
		return nil, err
	}
	return &order.OrderRes{
		List: list,
		Total: total,
	}, nil
}