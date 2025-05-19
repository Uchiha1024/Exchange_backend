package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"grpc-common/ucenter/types/withdraw"
	"mscoin-common/msdb"
	"mscoin-common/msdb/tran"
	"mscoin-common/op"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"
	"ucenter/internal/model"
	"ucenter/internal/svc"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	memberAddressDomain *domain.MemberAddressDomain
	memberDomain        *domain.MemberDomain
	memberWalletDomain  *domain.MemberWalletDomain
	transaction         tran.Transaction
	withdrawDomain      *domain.WithdrawDomain
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		memberAddressDomain: domain.NewMemberAddressDomain(svcCtx.Db),
		memberDomain:        domain.NewMemberDomain(svcCtx.Db),
		transaction:         tran.NewTransaction(svcCtx.Db.Conn),
		memberWalletDomain:  domain.NewMemberWalletDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.Cache),
		withdrawDomain:      domain.NewWithdrawDomain(svcCtx.Db, svcCtx.MarketRpc, svcCtx.BitcoinAddress),
	}
}

func (l *WithdrawLogic) FindAddressByCoinId(req *withdraw.WithdrawReq) (*withdraw.AddressSimpleList, error) {
	list, err := l.memberAddressDomain.FindAddressList(l.ctx, req.UserId, req.CoinId)
	if err != nil {
		return nil, err
	}
	var addressList []*withdraw.AddressSimple
	copier.Copy(&addressList, list)
	return &withdraw.AddressSimpleList{
		List: addressList,
	}, nil
}

func (l *WithdrawLogic) SendCode(req *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	//假设发送了一条短信 验证码是123456
	code := "123456"
	// 将验证码写入缓存
	err := l.svcCtx.Cache.SetWithExpireCtx(l.ctx, "WITHDRAW::"+req.Phone, code, time.Minute*10)
	return &withdraw.NoRes{}, err

}

func (l *WithdrawLogic) WithdrawCode(req *withdraw.WithdrawReq) (*withdraw.NoRes, error) {
	// 1. 从缓存中获取验证码
	userInfo, err := l.memberDomain.FindMemberById(l.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var redisCode string
	err = l.svcCtx.Cache.GetCtx(l.ctx, "WITHDRAW::"+userInfo.MobilePhone, &redisCode)
	if err != nil {
		return nil, err
	}
	if redisCode != req.Code {
		return nil, errors.New("验证码错误")
	}
	//校验密码
	if userInfo.JyPassword != req.JyPassword {
		return nil, errors.New("密码错误")
	}
	// 3. 根据用户id 和 unit 查询用户的钱包 判断余额是否足够
	memberWallet, err := l.memberWalletDomain.FindWalletByMemIdAndCoin(l.ctx, req.UserId, req.Unit)
	if err != nil {
		return nil, err
	}
	if memberWallet == nil {
		return nil, errors.New("钱包不存在")
	}
	if memberWallet.Balance < req.Amount {
		return nil, errors.New("余额不足")
	}
	err = l.transaction.Action(func(conn msdb.DbConn) error {
		//事务处理
		//4. 冻结用户的钱 提现币 经过比特币网络 需要时间
		err2 := l.memberWalletDomain.Freeze(l.ctx, conn, req.UserId, req.Amount, req.Unit)
		if err2 != nil {
			return err2
		}
		//5. 记录用户的提现
		wr := &model.WithdrawRecord{}
		wr.CoinId = memberWallet.CoinId
		wr.Address = req.Address
		wr.Fee = req.Fee
		wr.TotalAmount = req.Amount
		wr.ArrivedAmount = op.SubFloor(req.Amount, req.Fee, 10)
		wr.Remark = ""
		wr.CanAutoWithdraw = 0
		wr.IsAuto = 0
		wr.Status = 0 //审核中
		wr.CreateTime = time.Now().UnixMilli()
		wr.DealTime = 0
		wr.MemberId = req.UserId
		wr.TransactionNumber = "" //目前还没有交易编号
		var err error
		err = l.withdrawDomain.SaveRecord(l.ctx, wr)
		if err != nil {
			return err
		}
		// 6. 发送用户的提现事件到MQ当中 MQ消费者去处理提现（创建交易 广播到比特币的网络）
		marshaled, _ := json.Marshal(wr)
		data := database.KafkaData{
			Topic: "withdraw",
			Data:  marshaled,
			Key:   []byte(fmt.Sprintf("%d", req.UserId)),
		}
		for i := 0; i < 3; i++ {
			err = l.svcCtx.KafkaCli.SendSync(data)
			if err != nil {
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &withdraw.NoRes{}, nil
}


func (l *WithdrawLogic) WithdrawRecord(req *withdraw.WithdrawReq) (*withdraw.RecordList, error) {
	list, total, err := l.withdrawDomain.RecordList(l.ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	var rList []*withdraw.WithdrawRecord
	copier.Copy(&rList, list)
	return &withdraw.RecordList{
		List:  rList,
		Total: total,
	}, nil
}