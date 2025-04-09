package logic

import (
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/domain"
	"mscoin-common/tools"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OkxConfig struct {
	ApiKey    string
	SecretKey string
	Pass      string
	Host      string
	Proxy     string
}

type OkxResult struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

type Kline struct {
	wg          sync.WaitGroup
	okx         OkxConfig
	klineDomain *domain.KlineDomain
}

func NewKline(okx OkxConfig, mongoClient *database.MongoClient) *Kline {
	return &Kline{
		okx:         okx,
		klineDomain: domain.NewKlineDomain(mongoClient),
	}
}

func (k *Kline) Do(period string) {
	k.wg.Add(2)
	go k.getKlineData("BTC-USDT", "BTC/USDT", period)
	go k.getKlineData("ETH-USDT", "ETH/USDT", period)
	k.wg.Wait()

}

func (k *Kline) getKlineData(instId string, symbol string, period string) {
	// 获取K线数据
	api := k.okx.Host + "/api/v5/market/candles?instId=" + instId + "&bar=" + period
	timestamp := tools.ISO(time.Now())
	sign := tools.ComputeHmacSha256(timestamp+"GET"+"/api/v5/market/candles?instId="+instId+"&bar="+period, k.okx.SecretKey)
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = k.okx.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = k.okx.Pass
	resp, err := tools.GetWithHeader(api, header, k.okx.Proxy)
	if err != nil {
		logx.Errorw("getKlineData", logx.Field("err", err))
		k.wg.Done()
		return
	}

	var result OkxResult
	err = json.Unmarshal(resp, &result)
	if err != nil {
		logx.Errorw("json unmarshal", logx.Field("err", err))
		k.wg.Done()
		return
	}

	logx.Info("==================执行存储mongo====================")

	if result.Code == "0" {
		err = k.klineDomain.SaveBatch(result.Data, symbol, period)
		if err != nil {
			logx.Errorw("save mongo failed", logx.Field("err", err))
			k.wg.Done()
			return
		}
	
	}
	k.wg.Done()
	

	logx.Info("==================End====================")
}
