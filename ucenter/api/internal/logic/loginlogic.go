package logic

import (
	"context"

	"ucenter/api/internal/svc"
	"ucenter/api/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *login.LoginReq) (*login.LoginRes, error) {
	// todo: add your logic here and delete this line

	return &login.LoginRes{}, nil
}
