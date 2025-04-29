package logic

import (
	"context"
	"grpc-common/ucenter/types/member"
	"ucenter/internal/domain"
	"ucenter/internal/svc"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type MemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	memberDomain *domain.MemberDomain
}

func NewMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberLogic {
	return &MemberLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		memberDomain: domain.NewMemberDomain(svcCtx.Db),
	}
}

func (l *MemberLogic) FindMemberById(in *member.MemberReq) (*member.MemberInfo, error) {
	mem, err := l.memberDomain.FindMemberById(l.ctx, in.MemberId)
	if err != nil {
		logx.Errorw("MemberLogic-FindMemberById-ERROR", logx.Field("err", err))
		return nil, err
	}
	resp := &member.MemberInfo{}
	err = copier.Copy(resp, mem)
	if err != nil {
		logx.Errorw("MemberLogic-FindMemberById-Copy-ERROR", logx.Field("err", err))
		return nil, err
	}
	return resp, nil

}
