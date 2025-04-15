package logic

import (
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/domain"
	"mscoin-common/tools"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
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
	queueDomain *domain.QueueDomain
	redisCache  cache.Cache
}

func NewKline(okx OkxConfig, mongoClient *database.MongoClient,kafkaCli *database.KafkaClient, cache2 cache.Cache) *Kline {
	return &Kline{
		okx:         okx,
		klineDomain: domain.NewKlineDomain(mongoClient),
		queueDomain: domain.NewQueueDomain(kafkaCli),
		redisCache: cache2,

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
		
		if period == "1m" {
			//把这个最新的数据result.Data[0] 推送到market服务，推送到前端页面，实时进行变化
		//->kafka->market kafka消费者进行数据消费-> 通过websocket通道发送给前端 ->前端更新数据
			if len(result.Data)>0 {
				kafkaData := result.Data[0]
				k.queueDomain.Send1mKline(kafkaData,symbol)
				// 存入 redis 保存最新价格
				redisKey := strings.ReplaceAll(instId,"-", "::")
				k.redisCache.Set(redisKey+"::RATE", kafkaData[4])


			}

		}


	
	}
	k.wg.Done()
	

	logx.Info("==================End====================")
}
