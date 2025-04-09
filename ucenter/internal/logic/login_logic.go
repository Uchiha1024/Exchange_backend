package logic

import (
	"context"
	"errors"
	"mscoin-common/tools"
	"time"

	"grpc-common/ucenter/types/login"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *LoginLogic) Login(in *login.LoginReq) (*login.LoginRes, error) {

	// 先验证人机是否通过
	/* 	isVerify := l.CaptchaDomain.VerifyCaptcha(
	   		in.Captcha.Server,
	   		l.svcCtx.Config.Captcha.Vid,
	   		l.svcCtx.Config.Captcha.Key,
	   		in.Captcha.Token,
	   		2,
	   		in.Ip,
	   	)
	   	if !isVerify {
	   		return nil, errors.New("人机验证失败")
	   	}
	*/
	logx.Info("人机验证通过")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 验证用户名密码
	member, err := l.MemberDomain.FindByPhone(ctx, in.GetUsername())
	if err != nil {
		logx.Errorf("find by phone error: %v", err)
		return nil, errors.New("login failed")
	}

	if member == nil {
		return nil, errors.New("user not registered")
	}

	password := member.Password
	salt := member.Salt
	verifyResult := tools.Verify(in.GetPassword(), salt, password, nil)
	if !verifyResult {
		return nil, errors.New("password is incorrect")
	}

	// 3. 登录成功，生成token，提供给前端，前端调用传递token，我们进行token认证即可
	key := l.svcCtx.Config.JWT.AccessSecret
	expire := l.svcCtx.Config.JWT.AccessExpire
	token, err := tools.GetJwtToken(key, time.Now().Unix(), expire, member.Id)
	if err != nil {
		logx.Errorf("get jwt token error: %v", err)
		return nil, errors.New("token generate failed")
	}
	// 返回登录所需信息
	loginCount := member.LoginCount + 1
	go func(userId int64) {
		l.MemberDomain.UpdateLoginCount(context.Background(), userId, 1)
	}(member.Id)


	return &login.LoginRes{
		Token:         token,
		Id:            member.Id,
		Username:      member.Username,
		MemberLevel:   member.MemberLevelStr(),
		MemberRate:    member.MemberRate(),
		RealName:      member.RealName,
		Country:       member.Country,
		Avatar:        member.Avatar,
		PromotionCode: member.PromotionCode,
		SuperPartner:  member.SuperPartner,
		LoginCount:    int32(loginCount),
	}, nil
}
