package domain

import (
	"context"
	"errors"
	"mscoin-common/msdb"
	"mscoin-common/tools"
	"regexp"
	"ucenter/internal/dao"
	"ucenter/internal/model"
	"ucenter/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberDomain struct {
	MemberRepo repo.MemberRepo
}



func NewMemberDomain(db *msdb.MsDB) *MemberDomain {
	return &MemberDomain{
		MemberRepo: dao.NewMemberDao(db),
	}
}


func (m *MemberDomain) UpdateLoginCount(ctx context.Context, id int64, step int) {
	err := m.MemberRepo.UpdateLoginCount(ctx, id, step)
	if err != nil {
		logx.Errorf("update login count error: %v", err)
	}
}


func (m *MemberDomain) Register(
	ctx context.Context,
	phone string,
	password string,
	username string,
	country string,
	partner string,
	promotion string) error {

	member := model.NewMember()
	//member表字段比较多，所有的字段都不为null，也就是很多字段要填写默认值 写一个工具类 通过反射填充默认值即可
	_ = tools.Default(member)
	salt, pwd := tools.Encode(password, nil)
	member.Username = username
	member.Country = country
	member.Password = pwd
	member.MobilePhone = phone
	member.FillSuperPartner(partner)
	member.PromotionCode = promotion
	member.MemberLevel = model.GENERAL
	member.Salt = salt
	member.Avatar = "https://mszlu.oss-cn-beijing.aliyuncs.com/mscoin/defaultavatar.png"
	err := m.MemberRepo.Save(ctx, member)
	if err != nil {
		logx.Errorf("save member error: %v", err)
		return errors.New("database error")
	}
	return nil

}

func (m *MemberDomain) FindByPhone(context context.Context, phone string) (*model.Member, error) {
	// 1. 判断手机号是否为空
	if phone == "" {
		return nil, errors.New("phone is empty")
	}

	// 2. 判断手机号是否为11位
	if len(phone) != 11 {
		return nil, errors.New("phone is not valid")
	}

	// 3. 判断手机号是否为数字
	matched, err := regexp.MatchString("^[0-9]+$", phone)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("phone is not valid")
	}

	// 4. 数据库查询
	mem, err := m.MemberRepo.FindByPhone(context, phone)

	if err != nil {
		logx.Errorf("find by phone error: %v", err)
		return nil, errors.New("database error")
	}
	return mem, nil

}
