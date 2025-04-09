package domain

import "strings"

type ExchangeRateDomain struct {
}

func NewExchangeRateDomain() *ExchangeRateDomain {
	return &ExchangeRateDomain{}
}

func (d *ExchangeRateDomain) UsdRate(unit string) float64 {
	//应该据redis查询，在定时任务做一个根据实际的汇率接口 定期存入redis
	unit = strings.ToUpper(unit)
	if unit == "CNY" {
		return 6.95
	} else if unit == "JPY" {
		return 136.23
	}
	return 0

}
