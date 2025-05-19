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

func (s *WithdrawServer) SendCode(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	l := logic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.SendCode(in)
}

func (s *WithdrawServer) WithdrawCode(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	l := logic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.WithdrawCode(in)
}


func (s *WithdrawServer) WithdrawRecord(ctx context.Context, in *withdraw.WithdrawReq) (*withdraw.RecordList, error) {
	l := logic.NewWithdrawLogic(ctx, s.svcCtx)
	return l.WithdrawRecord(in)
}