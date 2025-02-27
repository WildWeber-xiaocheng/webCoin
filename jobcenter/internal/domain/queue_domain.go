package domain

import (
	"encoding/json"
	"jobcenter/internal/database"
	"jobcenter/internal/model"
)

type QueueDomain struct {
	kafkaCli *database.KafkaClient
}

const KLINE1M = "kline_1m"

// Send1mKline 向kafka发送最近1m的k线数据
func (d *QueueDomain) Send1mKline(data []string, symbol, period string) {
	kline := model.NewKline(data, period)
	bytes, _ := json.Marshal(kline)
	sendData := database.KafkaData{
		Topic: KLINE1M,
		Key:   []byte(symbol),
		Data:  bytes,
	}
	d.kafkaCli.Send(sendData)
}

func NewQueueDomain(kafkaCli *database.KafkaClient) *QueueDomain {
	return &QueueDomain{
		kafkaCli: kafkaCli,
	}
}
