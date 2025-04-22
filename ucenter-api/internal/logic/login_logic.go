package logic

import (
	"context"
	"grpc-common/ucenter/types/login"
	"mscoin-common/tools"
	"time"

	"ucenter-api/internal/svc"
	"ucenter-api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}



func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


func (l *LoginLogic) CheckLogin(token string) (bool, error) {
	_, err := tools.ParseToken(token, l.svcCtx.Config.JWT.AccessSecret)
	if err != nil {
		logx.Errorf("check login error: %v", err)
		return false, nil
	}
	return true, nil
}



func (l *LoginLogic) Login(req *types.LoginReq) ( *types.LoginRes, error) {
	// todo: add your logic here and delete this line

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	loginReq := &login.LoginReq{}
	if err := copier.Copy(loginReq, req); err != nil {
		return nil, err
	}

	loginResp, err := l.svcCtx.UCLoginRpc.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	resp := &types.LoginRes{}
	if err := copier.Copy(resp, loginResp); err != nil {
		return nil, err
	}

	return resp, err
}
