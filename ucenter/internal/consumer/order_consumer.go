package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"grpc-common/exchange/eclient"
	"grpc-common/exchange/types/order"
	"time"
	"ucenter/internal/database"
	"ucenter/internal/domain"
	"webCoin-common/enum"
	"webCoin-common/msdb"
	"webCoin-common/msdb/tran"
	"webCoin-common/op"
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
func ExchangeOrderAdd(redisClient *redis.Redis, client *database.KafkaClient, orderRpc eclient.Order, db *msdb.MsDB) {
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
		//使用go-zero自带的分布式锁
		lock := redis.NewRedisLock(redisClient, "exchange_order::"+fmt.Sprintf("%d::%s", addData.UserId, orderId))
		//查询订单信息 如果是正在交易中 继续 否则return
		acquireCtx, err := lock.AcquireCtx(ctx)
		if err != nil {
			logx.Error(err)
			logx.Info("已经有别的进程处理此消息")
			continue
		}
		if acquireCtx {
			transaction := tran.NewTransaction(db.Conn)
			walletDomain := domain.NewMemberWalletDomain(db, nil)
			err = transaction.Action(func(conn msdb.DbConn) error {
				if addData.Direction == 0 { //买
					//buy baseSymbol
					err := walletDomain.Freeze(ctx, conn, addData.UserId, addData.Money, addData.BaseSymbol)
					return err
				} else if addData.Direction == 1 { //卖
					//sell coinSymbol
					err := walletDomain.Freeze(ctx, conn, addData.UserId, addData.Money, addData.CoinSymbol)
					return err
				}
				return nil
			})
			if err != nil {
				logx.Error(err)
				cancelOrder(client, orderRpc, ctx, orderId, exchangeOrder, kafkaData)
				continue
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
				logx.Info("发送exchange_order_init_complete 消息成功：" + orderId)
				break
			}
			lock.Release()
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

type ExchangeOrder struct {
	Id            int64   `gorm:"column:id" json:"id"`
	OrderId       string  `gorm:"column:order_id" json:"orderId"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
	BaseSymbol    string  `gorm:"column:base_symbol" json:"baseSymbol"`
	CanceledTime  int64   `gorm:"column:canceled_time" json:"canceledTime"`
	CoinSymbol    string  `gorm:"column:coin_symbol" json:"coinSymbol"`
	CompletedTime int64   `gorm:"column:completed_time" json:"completedTime"`
	Direction     int     `gorm:"column:direction" json:"direction"`
	MemberId      int64   `gorm:"column:member_id" json:"memberId"`
	Price         float64 `gorm:"column:price" json:"price"`
	Status        int     `gorm:"column:status" json:"status"`
	Symbol        string  `gorm:"column:symbol" json:"symbol"`
	Time          int64   `gorm:"column:time" json:"time"`
	TradedAmount  float64 `gorm:"column:traded_amount" json:"tradedAmount"`
	Turnover      float64 `gorm:"column:turnover" json:"turnover"`
	Type          int     `gorm:"column:type" json:"type"`
	UseDiscount   string  `gorm:"column:use_discount" json:"useDiscount"`
}

// status
const (
	Trading = iota
	Completed
	Canceled
	OverTimed
	Init
)

var StatusMap = enum.Enum{
	Trading:   "TRADING",
	Completed: "COMPLETED",
	Canceled:  "CANCELED",
	OverTimed: "OVERTIMED",
}

// direction
const (
	BUY = iota
	SELL
)

var DirectionMap = enum.Enum{
	BUY:  "BUY",
	SELL: "SELL",
}

// type
const (
	MarketPrice = iota
	LimitPrice
)

var TypeMap = enum.Enum{
	MarketPrice: "MARKET_PRICE",
	LimitPrice:  "LIMIT_PRICE",
}

// 接收kafka的消息（从exchange/internal/consumer/trade_consumer.go中的readOrderComplete函数来的消息）
// 接收到的消息的是订单已完成，这里要修改用户的钱包信息
func ExchangeOrderComplete(redisCli *redis.Redis, cli *database.KafkaClient, db *msdb.MsDB) {
	//先接收消息
	for {
		kafkaData := cli.Read()
		var order *ExchangeOrder
		json.Unmarshal(kafkaData.Data, &order)
		if order == nil {
			continue
		}
		if order.Status != Completed {
			continue
		}
		logx.Info("收到exchange_order_complete_update_success 消息成功:" + order.OrderId)
		walletDomain := domain.NewMemberWalletDomain(db, nil)
		//基于redis的分布式锁
		lock := redis.NewRedisLock(redisCli, fmt.Sprintf("order_complete_update_wallet::%d", order.MemberId))
		acquire, err := lock.Acquire()
		if err != nil {
			logx.Error(err)
			logx.Info("有进程已经拿到锁进行处理了")
			continue
		}
		if acquire {
			// BTC/USDT
			ctx := context.Background()
			if order.Direction == BUY {
				//买  用baseSymbol
				baseWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.BaseSymbol)
				if err != nil {
					logx.Error(err)
					//重新将数据放入kafka
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				coinWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.CoinSymbol)
				if err != nil {
					logx.Error(err)
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				if order.Type == MarketPrice {
					//市价买的情况，则amount（以USDT为单位）为冻结的钱  order.turnover是成交额的钱，则需要还回去的钱为amount-order.turnover
					//冻结的钱进行解冻
					baseWallet.FrozenBalance = op.SubFloor(baseWallet.FrozenBalance, order.Amount, 8)
					//剩余没花完的钱还回去
					baseWallet.Balance = op.AddFloor(baseWallet.Balance, op.SubFloor(order.Amount, order.Turnover, 8), 8)
					coinWallet.Balance = op.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
				} else {
					//限价买的情况，冻结的钱是 order.price*amount  成交了turnover 还回去的钱为order.price*amount-order.turnover
					floor := op.MulFloor(order.Price, order.Amount, 8)
					//解冻
					baseWallet.FrozenBalance = op.SubFloor(baseWallet.FrozenBalance, floor, 8)
					//剩余没花完的钱还回去
					baseWallet.Balance = op.AddFloor(baseWallet.Balance, op.SubFloor(floor, order.Turnover, 8), 8)
					coinWallet.Balance = op.AddFloor(coinWallet.Balance, order.TradedAmount, 8)
				}
				err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
				if err != nil {
					logx.Error(err)
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
			} else { //卖
				//卖 不管是市价还是限价 都是卖的 BTC  解冻amount 得到的钱是 order.turnover
				coinWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.CoinSymbol)
				if err != nil {
					logx.Error(err)
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				baseWallet, err := walletDomain.FindWalletByMemIdAndCoin(ctx, order.MemberId, order.BaseSymbol)
				if err != nil {
					logx.Error(err)
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
				//解冻
				coinWallet.FrozenBalance = op.SubFloor(coinWallet.FrozenBalance, order.Amount, 8)
				//将成交额加上
				baseWallet.Balance = op.AddFloor(baseWallet.Balance, order.Turnover, 8)
				err = walletDomain.UpdateWalletCoinAndBase(ctx, baseWallet, coinWallet)
				if err != nil {
					logx.Error(err)
					cli.RPut(kafkaData)
					time.Sleep(250 * time.Millisecond)
					lock.Release()
					continue
				}
			}
			logx.Info("更新钱包成功:" + order.OrderId)
			lock.Release()
		}

	}
}
