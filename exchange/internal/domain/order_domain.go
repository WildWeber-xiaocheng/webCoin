package domain

import (
	"context"
	"errors"
	"exchange/internal/dao"
	"exchange/internal/model"
	"exchange/internal/repo"
	"fmt"
	"grpc-common/market/mclient"
	"grpc-common/ucenter/ucclient"
	"time"
	"webCoin-common/msdb"
	"webCoin-common/op"
	"webCoin-common/tools"
)

type ExchangeOrderDomain struct {
	orderRepo repo.ExchangeOrderRepo
}

func (d *ExchangeOrderDomain) FindOrderHistory(
	ctx context.Context,
	symbol string,
	page int64,
	size int64,
	memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderHistory(ctx, symbol, page, size, memberId)
	if err != nil {
		return nil, 0, err
	}
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func (d *ExchangeOrderDomain) FindOrderCurrent(
	ctx context.Context,
	symbol string,
	page int64,
	size int64,
	memberId int64) ([]*model.ExchangeOrderVo, int64, error) {
	list, total, err := d.orderRepo.FindOrderCurrent(ctx, symbol, page, size, memberId)
	if err != nil {
		return nil, 0, err
	}
	lv := make([]*model.ExchangeOrderVo, len(list))
	for i, v := range list {
		lv[i] = v.ToVo()
	}
	return lv, total, err
}

func (d *ExchangeOrderDomain) FindCurrentTradingCount(ctx context.Context, userId int64, symbol string, direction string) (int64, error) {
	return d.orderRepo.FindCurrentTradingCount(ctx, userId, symbol, model.DirectionMap.Code(direction))
}

// AddOrder 添加订单
// return float64:需要冻结的金额
func (d *ExchangeOrderDomain) AddOrder(
	ctx context.Context, conn msdb.DbConn,
	order *model.ExchangeOrder, coin *mclient.ExchangeCoin,
	baseWallet *ucclient.MemberWallet, coinWallet *ucclient.MemberWallet) (float64, error) {
	//设置订单状态为正在交易，交易数量初始化为0，订单创建时间为当前时间
	order.Status = model.Trading
	order.TradedAmount = 0
	order.Time = time.Now().UnixMilli()
	//todo 使用雪花算法
	order.OrderId = tools.Unq("E")
	var money float64
	//交易时，coin.Fee是费率，即手续费，这里实现的时候不考虑手续费问题
	//买 花USDT 卖 用BTC
	//买的时候如果是市价，则冻结的直接就是amount
	//todo 这里按照文档来写的，和视频P41不一致
	if order.Direction == model.BUY {
		var turnover float64 = 0
		if order.Type == model.MarketPrice {
			turnover = order.Amount
		} else {
			turnover = op.MulN(order.Amount, order.Price, 5)
		}
		//费率
		fee := op.MulN(turnover, coin.Fee, 5)
		if baseWallet.Balance < turnover {
			return 0, errors.New("余额不足")
		}
		if baseWallet.Balance-turnover < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		//需要冻结的钱 turnover+fee
		money = op.AddN(turnover, fee, 5)
	} else {
		//卖
		fee := op.MulN(order.Amount, coin.Fee, 5)
		if coinWallet.Balance < order.Amount {
			return 0, errors.New("余额不足")
		}
		if coinWallet.Balance-order.Amount < fee {
			return 0, errors.New("手续费不足 需要:" + fmt.Sprintf("%f", fee))
		}
		money = op.AddN(order.Amount, fee, 5)
	}
	err := d.orderRepo.Save(ctx, conn, order)
	if err != nil {
		return 0, err
	}
	return money, nil
}

func (d *ExchangeOrderDomain) FindByOrderId(ctx context.Context, orderId string) (*model.ExchangeOrder, error) {
	exchangeOrder, err := d.orderRepo.FindByOrderId(ctx, orderId)
	if err == nil && exchangeOrder == nil {
		return nil, errors.New("订单号不存在")
	}
	return exchangeOrder, err
}

func (d *ExchangeOrderDomain) UpdateOrderStatusCancel(ctx context.Context, orderId string, updateStatus int) error {
	//todo 这里是按照文档来写的，不是视频P44
	return d.orderRepo.UpdateOrderStatusCancel(ctx, orderId, model.Canceled, updateStatus, time.Now().UnixMilli())
}

func NewExchangeOrderDomain(db *msdb.MsDB) *ExchangeOrderDomain {
	return &ExchangeOrderDomain{orderRepo: dao.NewExchangeOrderDao(db)}
}
