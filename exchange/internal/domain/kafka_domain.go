package domain

import (
	"context"
	"encoding/json"
	"exchange/internal/database"
	"exchange/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type KafkaDomain struct {
	kafkaClient *database.KafkaClient
	orderDomain *ExchangeOrderDomain
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

type AddOrderResult struct {
	UserId  int64  `json:"userId"`
	OrderId string `json:"orderId"`
}

func (d *KafkaDomain) WaitAddOrderResult(topic string) {
	cli := d.kafkaClient.StartRead(topic)
	for {
		kafkaData := cli.Read()
		logx.Info("收到订单增加结果:" + string(kafkaData.Data))
		var orderResult AddOrderResult
		json.Unmarshal(kafkaData.Data, &orderResult)
		exchangeOrder, err := d.orderDomain.orderRepo.FindByOrderId(context.Background(), orderResult.OrderId)
		if err != nil {
			logx.Error(err)
			err := d.orderDomain.UpdateOrderStatusCancel(context.Background(), orderResult.OrderId)
			if err != nil {
				logx.Error(err)
				d.kafkaClient.RPut(kafkaData)
			}
			continue
		}
		if exchangeOrder == nil {
			logx.Error("订单id不存在")
			continue
		}
		if exchangeOrder.Status != model.Init {
			logx.Error("订单已经被处理过")
			continue
		}
		err = d.orderDomain.UpdateOrderStatusTrading(context.Background(), orderResult.OrderId)
		if err != nil {
			logx.Error(err)
			d.kafkaClient.RPut(kafkaData)
			continue
		}

		//订单初始化完成 订单需要加入到撮合交易当中
		//发送消息到kafka 等待撮合交易引擎进行交易撮合 如果没有撮合交易成功，加入到撮合交易队列，继续等待完成撮合
		exchangeOrder.Status = model.Trading
		for {
			bytes, _ := json.Marshal(orderResult)
			orderData := database.KafkaData{
				Topic: "exchange_order_trading",
				Key:   []byte(orderResult.OrderId),
				Data:  bytes,
			}
			err := d.kafkaClient.SendSync(orderData)
			if err != nil { //等待一段时间再次发送
				logx.Error(err)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			logx.Info("订单创建成功，发送创建成功消息:", orderResult.OrderId)
			break
		}
	}
}

func NewKafkaDomain(kafkaClient *database.KafkaClient, orderDomain *ExchangeOrderDomain) *KafkaDomain {
	k := &KafkaDomain{
		kafkaClient: kafkaClient,
		orderDomain: orderDomain,
	}
	go k.WaitAddOrderResult("exchange_order_init_complete")
	return k
}
