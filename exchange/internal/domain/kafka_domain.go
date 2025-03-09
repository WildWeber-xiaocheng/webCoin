package domain

import (
	"encoding/json"
	"exchange/internal/database"
	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaDomain struct {
	kafkaClient *database.KafkaClient
}

func (d *KafkaDomain) SendOrderAdd(
	topic string,
	userId int64,
	orderId string,
	money float64,
	symbol string,
	direction int,
	baseSymbol string,
	coinSymbol string) error {
	m := make(map[string]any)
	m["userId"] = userId
	m["orderId"] = orderId
	m["money"] = money
	m["symbol"] = symbol
	m["direction"] = direction
	m["baseSymbol"] = baseSymbol
	m["coinSymbol"] = coinSymbol
	marshal, _ := json.Marshal(m)
	data := database.KafkaData{
		Topic: topic,
		Key:   []byte(orderId),
		Data:  marshal,
	}
	logx.Info(string(marshal))
	return d.kafkaClient.SendSync(data)
}

func NewKafkaDomain(kafkaClient *database.KafkaClient) *KafkaDomain {
	return &KafkaDomain{
		kafkaClient: kafkaClient,
	}
}
