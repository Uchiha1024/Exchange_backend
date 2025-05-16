package server

import (
	"context"
	"grpc-common/ucenter/types/withdraw"
	"ucenter/internal/logic"
	"ucenter/internal/svc"
)

type WithdrawServer struct {
	svcCtx *svc.ServiceContext
	withdraw.UnimplementedWithdrawServer
}

func NewWithdrawServer(svcCtx *svc.ServiceContext) *WithdrawServer {
	return &WithdrawServer{
		svcCtx: svcCtx,
	}
}

func (s *WithdrawServer) FindAddressByCoinId(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.AddressSimpleList, error) {
	l := logic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.FindAddressByCoinId(in)
}
