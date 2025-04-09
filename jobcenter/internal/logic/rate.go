package logic

import (
	"encoding/json"
	"mscoin-common/tools"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type Rate struct {
	wg    sync.WaitGroup
	okx   OkxConfig
	Cache cache.Cache
}

func NewRate(okx OkxConfig, cache cache.Cache) *Rate {
	return &Rate{
		okx:   okx,
		Cache: cache,
	}
}

type OkxExchangeRateResult struct {
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
	Data []ExchangeRate `json:"data"`
}
type ExchangeRate struct {
	UsdCny string `json:"usdCny"`
}

var redisKey = "USDT::CNY::RATE"

func (r *Rate) Do() {
	r.wg.Add(1)
	go r.CnyUsdRate()
	r.wg.Wait()

}

func (r *Rate) CnyUsdRate() {
	// 通过请求okx的api获取cny/usd的汇率
	api := r.okx.Host + "/api/v5/market/exchange-rate"
	timestamp := tools.ISO(time.Now())
	sign := tools.ComputeHmacSha256(timestamp+"GET"+"/api/v5/market/exchange-rate", r.okx.SecretKey)
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = r.okx.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = r.okx.Pass
	res, err := tools.GetWithHeader(api, header, r.okx.Proxy)
	if err != nil {
		logx.Errorw("get cny/usd rate failed", logx.Field("error", err))
		r.wg.Done()
		return
	}

	var result OkxExchangeRateResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		logx.Errorw("unmarshal cny/usd rate failed", logx.Field("error", err))
		r.wg.Done()
		return
	}
	cnyUsdRate := result.Data[0].UsdCny
	// 将cny/usd的汇率写入redis
	err = r.Cache.Set(redisKey, cnyUsdRate)
	if err != nil {
		logx.Errorw("set cny/usd rate failed", logx.Field("error", err))
		r.wg.Done()
		return
	}
	r.wg.Done()

}
