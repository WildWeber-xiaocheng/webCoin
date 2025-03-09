package logic

import (
	"context"
	"errors"
	"exchange/internal/domain"
	"exchange/internal/model"
	"exchange/internal/svc"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/exchange/types/order"
	"grpc-common/market/types/market"
	"grpc-common/ucenter/types/member"
)

type ExchangeOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	exchangeOrderDomain *domain.ExchangeOrderDomain
}

func (l *ExchangeOrderLogic) FindOrderHistory(req *order.OrderReq) (*order.OrderRes, error) {
	exchangeOrders, total, err := l.exchangeOrderDomain.FindOrderHistory(
		l.ctx,
		req.Symbol,
		req.Page,
		req.PageSize,
		req.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

func (l *ExchangeOrderLogic) FindOrderCurrent(req *order.OrderReq) (*order.OrderRes, error) {
	exchangeOrders, total, err := l.exchangeOrderDomain.FindOrderCurrent(
		l.ctx,
		req.Symbol,
		req.Page,
		req.PageSize,
		req.UserId)
	if err != nil {
		return nil, err
	}
	var list []*order.ExchangeOrder
	copier.Copy(&list, exchangeOrders)
	return &order.OrderRes{
		List:  list,
		Total: total,
	}, nil
}

// Add 发布委托（添加订单）
func (l *ExchangeOrderLogic) Add(req *order.OrderReq) (*order.AddOrderRes, error) {
	//1.检查参数是否合法
	memberRes, err := l.svcCtx.MemberRpc.FindMemberById(l.ctx, &member.MemberReq{
		MemberId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if memberRes.TransactionStatus == 0 {
		return nil, errors.New("此用户已经被禁止交易")
	}
	if req.Type == model.TypeMap[model.LimitPrice] && req.Price <= 0 {
		return nil, errors.New("限价模式下价格不能小于等于0")
	}
	if req.Amount <= 0 {
		return nil, errors.New("数量不能小于等于0")
	}
	exchangeCoin, err := l.svcCtx.MarketRpc.FindSymbolInfo(l.ctx, &market.MarketReq{
		Symbol: req.Symbol,
	})
	if err != nil {
		return nil, errors.New("nonsupport coin")
	}
	if exchangeCoin.Exchangeable != 1 && exchangeCoin.Enable != 1 {
		return nil, errors.New("coin forbidden")
	}
	//基准币 例如symbol是BTC/USDT 则基准币是USDT
	baseSymbol := exchangeCoin.GetBaseSymbol()
	//交易币
	coinSymbol := exchangeCoin.GetCoinSymbol()
	cc := baseSymbol
	//卖则根据交易币来查 ，买则根据基准币来查
	if req.Direction == model.DirectionMap[model.SELL] {
		//根据交易币查询
		cc = coinSymbol
	}
	coin, err := l.svcCtx.MarketRpc.FindCoinInfo(l.ctx, &market.MarketReq{
		Unit: cc,
	})
	if err != nil || coin == nil {
		return nil, errors.New("nonsupport coin")
	}
	if req.Type == model.TypeMap[model.MarketPrice] && req.Direction == model.DirectionMap[model.BUY] {
		if exchangeCoin.GetMinTurnover() > 0 && req.Amount < float64(exchangeCoin.GetMinTurnover()) {
			return nil, errors.New("成交额至少是" + fmt.Sprintf("%d", exchangeCoin.GetMinTurnover()))
		}
	} else {
		if exchangeCoin.GetMaxVolume() > 0 && exchangeCoin.GetMaxVolume() < req.Amount {
			return nil, errors.New("数量超出" + fmt.Sprintf("%f", exchangeCoin.GetMaxVolume()))
		}
		if exchangeCoin.GetMinVolume() > 0 && exchangeCoin.GetMinVolume() > req.Amount {
			return nil, errors.New("数量不能低于" + fmt.Sprintf("%f", exchangeCoin.GetMinVolume()))
		}
	}
}

func (l *ExchangeOrderLogic) FindByOrderId(req *order.OrderReq) (*order.ExchangeOrderOrigin, error) {
	return nil, nil
}

func (l *ExchangeOrderLogic) CancelOrder(req *order.OrderReq) (*order.CancelOrderRes, error) {
	return nil, nil
}

func NewExchangeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderLogic {
	return &ExchangeOrderLogic{
		ctx:                 ctx,
		svcCtx:              svcCtx,
		Logger:              logx.WithContext(ctx),
		exchangeOrderDomain: domain.NewExchangeOrderDomain(svcCtx.Db),
	}
}
