package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/eclient"
	"grpc-common/exchange/types/order"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"
	"webCoin-common/msdb"
)

type OrderAdd struct {
	UserId     int64   `json:"userId"`
	OrderId    string  `json:"orderId"`
	Money      float64 `json:"money"`
	Symbol     string  `json:"symbol"`
	Direction  int     `json:"direction"`
	BaseSymbol string  `json:"baseSymbol"`
	CoinSymbol string  `json:"coinSymbol"`
}

// 消费kafka中订单数据
func ExchangeOrderAdd(client *database.KafkaClient, orderRpc eclient.Order, db *msdb.MsDB) {
	for {
		kafkaData := client.Read()
		//if kafkaData == nil {
		//	continue
		//}
		var addData OrderAdd
		//反序列化，得到数据
		err := json.Unmarshal(kafkaData.Data, &addData)

		if err != nil {
			//不是这个消息 消息类型错误
			logx.Error(err)
			continue
		}
		logx.Info("读取到订单添加消息：", string(kafkaData.Data))
		//获取订单id
		orderId := string(kafkaData.Key)
		if addData.OrderId != orderId {
			logx.Error(errors.New("不合法的消息，订单号不匹配"))
			continue
		}
		//查询订单信息 如果是正在交易中 继续 否则return
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		exchangeOrder, err := orderRpc.FindByOrderId(ctx, &order.OrderReq{
			OrderId: orderId,
		})
		if err != nil {
			logx.Error(err)
			cancelOrder(client, orderRpc, ctx, orderId, exchangeOrder, kafkaData)
			continue
		} else if exchangeOrder == nil {
			logx.Error("orderId: " + orderId + "不存在")
			continue
		} else if exchangeOrder.Status != 4 { //4表示订单的状态为init
			logx.Error("orderId: " + orderId + "已经被操作过")
			continue
		}

		walletDomain := domain.NewMemberWalletDomain(db)
		if addData.Direction == 0 { //买
			//buy baseSymbol
			err := walletDomain.Freeze(ctx, addData.UserId, addData.Money, addData.BaseSymbol)
			if err != nil {
				//取消订单，资金冻结失败
				cancelOrder(client, orderRpc, ctx, orderId, exchangeOrder, kafkaData)
				continue
			}
		} else if addData.Direction == 1 { //卖
			//sell coinSymbol
			err := walletDomain.Freeze(ctx, addData.UserId, addData.Money, addData.CoinSymbol)
			if err != nil {
				//重新消费
				cancelOrder(client, orderRpc, ctx, orderId, exchangeOrder, kafkaData)
				continue
			}
		}
		//都完成后 通知订单进行状态变更 需要保证一定发送成功
		//将订单的状态由init改为trading
		for {
			m := make(map[string]any)
			m["userId"] = addData.UserId
			m["orderId"] = orderId
			marshal, _ := json.Marshal(m)
			data := database.KafkaData{
				Topic: "exchange_order_init_complete",
				Key:   []byte(orderId),
				Data:  marshal,
			}
			err := client.SendSync(data)
			if err != nil {
				logx.Error(err)
				time.Sleep(250 * time.Millisecond)
				continue
			}
			break
		}
	}
}

// 因为取消订单可能失败，所以将orderRpc.CancelOrder再进行一层封装
func cancelOrder(client *database.KafkaClient, orderRpc eclient.Order,
	ctx context.Context, orderId string, exchangeOrder *eclient.ExchangeOrderOrigin,
	kafkaData database.KafkaData) {
	var status int32 = 0
	if exchangeOrder != nil {
		status = exchangeOrder.Status
	}
	_, err := orderRpc.CancelOrder(ctx, &order.OrderReq{
		OrderId:      orderId,
		UpdateStatus: status,
	})
	if err != nil { //订单取消失败，重新放入kafka重新进行消费
		client.RPut(kafkaData)
	}
}
