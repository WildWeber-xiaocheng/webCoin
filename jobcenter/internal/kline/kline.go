package kline

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"
	"webCoin-common/tools"
)

//var secretKey = "84792342C53C9976E039C9385EF5FB62"

type Kline struct {
	wg sync.WaitGroup
	c  OkxConfig
}

func NewKline(c OkxConfig) *Kline {
	return &Kline{
		wg: sync.WaitGroup{},
		c:  c,
	}
}

type OkxResult struct {
	Code string     `json:"code"`
	Msg  string     `json:"msg"`
	Data [][]string `json:"data"`
}

type OkxConfig struct {
	ApiKey    string
	SecretKey string
	Pass      string
	Host      string
	Proxy     string
}

// k线数据拉取
// period 时间间隔
func (k *Kline) Do(period string) {
	log.Println("============启动k线数据拉取==============")
	k.wg.Add(2)
	go k.getKlineData("BTC-USDT", "BTC/USDT", period)
	go k.getKlineData("ETH-USDT", "ETH/USDT", period)
	k.wg.Wait()
	log.Println("===============k线数据拉取结束===============")
}

// 通过调用外部api获取k先数据
func (k *Kline) getKlineData(instId string, symbol, period string) {
	api := "GET/api/v5/market/candles?instId=" + instId + "&bar=" + period
	timestamp := tools.ISO(time.Now())
	secretKey := k.c.SecretKey
	sha256 := tools.ComputeHmacSha256(timestamp+api, secretKey)
	sign := base64.StdEncoding.EncodeToString([]byte(sha256))
	header := make(map[string]string)
	header["OK-ACCESS-KEY"] = k.c.ApiKey
	header["OK-ACCESS-SIGN"] = sign
	header["OK-ACCESS-TIMESTAMP"] = timestamp
	header["OK-ACCESS-PASSPHRASE"] = k.c.Pass
	//系统的代理ip：http://127.0.0.1:7890
	respBody, err := tools.GetWithHeader(
		k.c.Host+"/api/v5/market/candles?instId="+instId+"&bar="+period,
		header,
		k.c.Proxy)
	if err != nil {
		log.Println(err)
		k.wg.Done()
		return
	}
	//log.Println(instId, string(respBody))
	result := &OkxResult{}
	err = json.Unmarshal(respBody, result)
	if err != nil {
		log.Println(err)
		k.wg.Done()
		return
	}
	log.Println("===========获取到k线数据=============")
	log.Println("instId:", instId, "period:", period)
	log.Println("result data:", string(respBody))
	log.Println("===========end=============")
}
