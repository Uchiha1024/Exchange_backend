package domain

import (
	"context"
	"jobcenter/internal/dao"
	"jobcenter/internal/database"
	"jobcenter/internal/model"
	"jobcenter/internal/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type KlineDomain struct {
	klineRepo repo.KlineRepo
}

func NewKlineDomain(client *database.MongoClient) *KlineDomain {
	return &KlineDomain{
		klineRepo: dao.NewKlineDao(client.Db),
	}
}

func (k *KlineDomain) SaveBatch(data [][]string, symbol, period string) error {
	klines := make([]*model.Kline, len(data))
	for i, v := range data {
		klines[i] = model.NewKline(v, period)

	}
	err := k.klineRepo.DeleteGtTime(context.Background(), klines[len(data)-1].Time, symbol, period)
	if err != nil {
		logx.Errorw("删除数据失败", logx.Field("err", err))
		return err

	}
	err = k.klineRepo.SaveBatch(context.Background(), klines, symbol, period)
	if err != nil {
		logx.Errorw("保存数据失败", logx.Field("err", err))
		return err
	}
	return nil
}
