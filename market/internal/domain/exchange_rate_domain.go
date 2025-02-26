package domain

import (
	"strings"
)

const CNY = "CNY"
const JPY = "JPY"
const HKD = "HKD"

type ExchangeRateDomain struct {
}

func NewExchangeRateDomain() *ExchangeRateDomain {
	return &ExchangeRateDomain{}
}

// 参数：unit:单位
// 返回值：汇率
func (d *ExchangeRateDomain) GetUsdRate(unit string) float64 {
	//应该根据redis查询实时汇率，在定时任务做一个根据实际的汇率接口 定期存入redis
	upper := strings.ToUpper(unit)
	if upper == CNY {
		return 7.00
	} else if upper == JPY {
		return 110.02
	} else if upper == HKD {
		return 7.8497
	}
	return 0
}
