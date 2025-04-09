package logic

import (
	"context"
	"grpc-common/ucenter/types/register"
	"time"

	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	logx.Info("api register")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	regReq := &register.RegReq{}
	if err := copier.Copy(regReq, req); err != nil {
		return nil, err
	}

	// 打印入参
	logx.Infof("接收到请求参数: %+v", regReq)
	_, err = l.svcCtx.UCRegisterRpc.RegisterByPhone(ctx, regReq)
	if err != nil {
		return nil, err
	}

	return
}

func (l *RegisterLogic) SendCode(req *types.SendCodeRequest) (resp *types.SendCodeResponse, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = l.svcCtx.UCRegisterRpc.SendCode(ctx, &register.CodeReq{
		Phone:   req.Phone,
		Country: req.Country,
	})
	if err != nil {
		return nil, err
	}

	return
}
