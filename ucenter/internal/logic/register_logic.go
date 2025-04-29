package logic

import (
	"context"
	"errors"
	"mscoin-common/tools"
	"time"

	"grpc-common/ucenter/types/register"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

const RegisterCacheKey = "REGISTER::"

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	CaptchaDomain *domain.CaptchaDomain
	MemberDomain  *domain.MemberDomain
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:           ctx,
		svcCtx:        svcCtx,
		Logger:        logx.WithContext(ctx),
		CaptchaDomain: domain.NewCaptchaDomain(),
		MemberDomain:  domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *RegisterLogic) RegisterByPhone(in *register.RegReq) (*register.RegRes, error) {
	// todo: add your logic here and delete this line
	logx.Info("ucenter rpc register by phone")

	// 参数验证
	// 打印入参
	logx.Infof("接收到请求参数: %+v", in)
	if in == nil {
		return nil, errors.New("请求参数不能为空")
	}

	if in.Phone == "" {
		logx.Error("手机号为空")
		return nil, errors.New("手机号不能为空")
	}

	if in.Code == "" {
		return nil, errors.New("验证码不能为空")
	}

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

	// 校验验证码 从 redis中 获取
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var redisCode string

	logx.Infof("获取redisKey: %s", RegisterCacheKey+in.Phone)

	// 获取原始 Redis 客户端

	err := l.svcCtx.Cache.GetCtx(ctx, RegisterCacheKey+in.Phone, &redisCode)
	if err != nil {
		logx.Errorf("获取验证码失败: %v", err)
		return nil, errors.New("获取验证码失败")
	}

	if redisCode != in.Code {
		return nil, errors.New("验证码错误")
	}

	mem, err := l.MemberDomain.FindByPhone(ctx, in.Phone)
	if err != nil {
		return nil, errors.New("server error,contact admin")
	}

	if mem != nil {
		return nil, errors.New("phone already registered")
	}

	// 生成member模型，存入数据库
	err = l.MemberDomain.Register(
		ctx,
		in.Phone,
		in.Password,
		in.Username,
		in.Country,
		in.SuperPartner,
		in.Promotion,
	)

	if err != nil {
		return nil, errors.New("register failed")
	}

	return &register.RegRes{}, nil
}

func (l *RegisterLogic) SendCode(in *register.CodeReq) (*register.NoRes, error) {
	code := tools.Rand4Code()

	// 假设调用短信平台
	go func() {
		logx.Info("调用短信平台")
	}()

	logx.Infof("发送验证码: %s", code)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := l.svcCtx.Cache.SetWithExpireCtx(ctx, RegisterCacheKey+in.Phone, code, 15*time.Minute)
	if err != nil {
		logx.Errorf("设置验证码缓存失败: %v", err)
		return nil, errors.New("设置验证码缓存失败")
	}

	return &register.NoRes{}, nil

}
