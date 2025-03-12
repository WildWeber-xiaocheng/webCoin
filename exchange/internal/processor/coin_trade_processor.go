package processor

import (
	"context"
	"encoding/json"
	"exchange/internal/database"
	"exchange/internal/domain"
	"exchange/internal/model"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"sort"
	"sync"
	"time"
	"webCoin-common/msdb"
	"webCoin-common/op"
)

// CoinTradeFactory 工厂 专门生产对应symbol的交易引擎
// 交易引擎的数量与exchange_coin表相对应
type CoinTradeFactory struct {
	tradeMap map[string]*CoinTrade
	mux      sync.RWMutex
}

func NewCoinTradeFactory() *CoinTradeFactory {
	return &CoinTradeFactory{
		tradeMap: make(map[string]*CoinTrade),
	}
}

// CoinTradeFactory初始化的操作
// 查询所有的exchange_coin内容 循环创建交易引擎
func (c *CoinTradeFactory) Init(marketRpc mclient.Market, cli *database.KafkaClient, db *msdb.MsDB) {
	ctx := context.Background()
	//拿到exchange_coin所有数据
	exchangeCoinRes, err := marketRpc.FindExchangeCoinVisible(ctx, &market.MarketReq{})
	if err != nil {
		logx.Error(err)
		return
	}
	for _, v := range exchangeCoinRes.List {
		c.AddCoinTrade(v.Symbol, NewCoinTrade(v.Symbol, cli, db))
	}
}

// 添加交易引擎
// symbol:币种 ct:交易引擎
func (c *CoinTradeFactory) AddCoinTrade(symbol string, ct *CoinTrade) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.tradeMap[symbol] = ct
}

func (c *CoinTradeFactory) GetCoinTrade(symbol string) *CoinTrade {
	//初始化的操作
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.tradeMap[symbol]
}

// 之所以这里定义一个结构，是因为要对这个数组进行排序
// 按照升序排序
type TradeTimeQueue []*model.ExchangeOrder

func (t TradeTimeQueue) Len() int {
	return len(t)
}
func (t TradeTimeQueue) Less(i, j int) bool {
	//升序
	return t[i].Time < t[j].Time
}
func (t TradeTimeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type LimitPriceQueue struct {
	mux  sync.RWMutex
	list TradeQueue
}

// 可以利用price来进行排序
type LimitPriceMap struct {
	price float64
	list  []*model.ExchangeOrder
}

// 降序的排序
type TradeQueue []*LimitPriceMap

func (t TradeQueue) Len() int {
	return len(t)
}
func (t TradeQueue) Less(i, j int) bool {
	//降序
	return t[i].price > t[j].price
}
func (t TradeQueue) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// CoinTrade 撮合交易引擎  每一个交易对 BTC/USDT 都有各自的一个引擎
type CoinTrade struct {
	symbol          string
	buyMarketQueue  TradeTimeQueue // 买 市价队列 存储未交易的订单
	bmMux           sync.RWMutex
	sellMarketQueue TradeTimeQueue
	smMux           sync.RWMutex
	buyLimitQueue   *LimitPriceQueue //买 限价队列 从高到低
	sellLimitQueue  *LimitPriceQueue //卖 限价队列 从低到高 因为从买家来看，肯定期望买入的价格越低越好，所以这里将卖的限价队列由低到高排序，这样就可以尽早卖出
	buyTradePlate   *TradePlate      //买盘
	sellTradePlate  *TradePlate      //卖盘
	kafkaClient     *database.KafkaClient
	db              *msdb.MsDB
}

// TradePlate 盘口信息 每一个Symbol对应一个盘口
type TradePlate struct {
	Items     []*TradePlateItem `json:"items"`
	Symbol    string
	direction int //sell/buy
	maxDepth  int //最多的Item的数量
	mux       sync.RWMutex
}

// 将订单添加到买卖盘中（添加到TradePlate的items中）
func (p *TradePlate) Add(order *model.ExchangeOrder) {
	if p.direction != order.Direction { //买/卖放错盘了
		return
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	if order.Type == model.MarketPrice {
		//市价不进入买卖盘
		//买卖盘 委托订单的形式下产生的一个金融说辞
		//委托 基本上可以认定是成交的，一旦发生委托  那么就意味着买方和卖方市场就已经形成了
		//买卖盘就是存放委托的地方，包含已成交、未成交、全部，这些都可以叫买卖盘
		//看我们的应用给用户展示哪方面的数据，我们选择展示未成交的买单和卖单
		return
	}
	size := len(p.Items)
	if size > 0 {
		for _, v := range p.Items {
			if v.Price == order.Price {
				//order.Amount= 20  order.TradedAmount=10 10
				v.Amount = op.FloorFloat(v.Amount+(order.Amount-order.TradedAmount), 8)
				return
			}
		}
	}
	if size < p.maxDepth {
		tpi := &TradePlateItem{
			Amount: op.FloorFloat(order.Amount-order.TradedAmount, 8),
			Price:  order.Price,
		}
		p.Items = append(p.Items, tpi)
	}
}

type TradePlateItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

func NewTradePlate(symbol string, direction int) *TradePlate {
	return &TradePlate{
		Symbol:    symbol,
		direction: direction,
		maxDepth:  100,
	}
}

type TradePlateResult struct {
	Direction    string            `json:"direction"`
	MaxAmount    float64           `json:"maxAmount"`
	MinAmount    float64           `json:"minAmount"`
	HighestPrice float64           `json:"highestPrice"`
	LowestPrice  float64           `json:"lowestPrice"`
	Symbol       string            `json:"symbol"`
	Items        []*TradePlateItem `json:"items"`
}

func (p *TradePlate) AllResult() *TradePlateResult {
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items
	return result
}

// 转换买卖盘数据格式为TradePlateResult，用于前端展示
// num:展示买卖盘中前num个item盘口数据
func (p *TradePlate) Result(num int) *TradePlateResult {
	if num > len(p.Items) {
		num = len(p.Items)
	}
	result := &TradePlateResult{}
	direction := model.DirectionMap.Value(p.direction)
	result.Direction = direction
	//这里的最高最低价格、最大最小数量都是遍历买卖盘的Items得出的
	result.MaxAmount = p.getMaxAmount()
	result.MinAmount = p.getMinAmount()
	result.HighestPrice = p.getHighestPrice()
	result.LowestPrice = p.getLowestPrice()
	result.Symbol = p.Symbol
	result.Items = p.Items[:num] //截取前num个item
	return result
}

func (p *TradePlate) getMaxAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = 0
	for _, v := range p.Items {
		if v.Amount > amount {
			amount = v.Amount
		}
	}
	return amount
}

func (p *TradePlate) getMinAmount() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var amount float64 = p.Items[0].Amount
	for _, v := range p.Items {
		if v.Amount < amount {
			amount = v.Amount
		}
	}
	return amount
}

func (p *TradePlate) getHighestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = 0
	for _, v := range p.Items {
		if v.Price > price {
			price = v.Price
		}
	}
	return price
}
func (p *TradePlate) getLowestPrice() float64 {
	if len(p.Items) <= 0 {
		return 0
	}
	var price float64 = p.Items[0].Price
	for _, v := range p.Items {
		if v.Price < price {
			price = v.Price
		}
	}
	return price
}

// 修改买卖盘中交易的订单的相关数据
func (p *TradePlate) Remove(order *model.ExchangeOrder, amount float64) {
	for i, v := range p.Items {
		if v.Price == order.Price {
			v.Amount = op.SubFloor(v.Amount, amount, 8)
			if v.Amount <= 0 { //从买卖盘中移出该订单
				p.Items = append(p.Items[:i], p.Items[i+1:]...)
			}
			break
		}
	}
}

// 撮合交易核心代码
func (t *CoinTrade) Trade(exchangeOrder *model.ExchangeOrder) {
	//1. 当订单进来之后，我们判断 buy还是sell
	//2. 确定 市价 还是限价
	//3. buy 就和 sell队列进行匹配
	//4. sell 就和买的队列进行匹配
	//5. exchangeOrder 还未交易的，放入买卖盘 同时放入 交易引擎的队列中 等待下次匹配
	//6. 订单就会更新 订单的状态要变 冻结的金额 扣除等等
	//if exchangeOrder.Direction == model.BUY {
	//	//放入买盘卖盘 要发送页面上显示  把结果发送到kafka当中
	//	t.buyTradePlate.Add(exchangeOrder)
	//	t.sendTradPlateMsg(t.buyTradePlate)
	//} else {
	//	t.sellTradePlate.Add(exchangeOrder)
	//	t.sendTradPlateMsg(t.sellTradePlate)
	//}

	//exchangeOrder 买 和卖的队列进行匹配 还是卖 和买的队列进行匹配
	//先进行一轮交易，对于未交易完成的订单进入买卖盘
	var limitPriceList *LimitPriceQueue
	var marketPriceList TradeTimeQueue
	if exchangeOrder.Direction == model.BUY { //如果是买，则需要和卖盘/卖的市价队列去匹配
		limitPriceList = t.sellLimitQueue
		marketPriceList = t.sellMarketQueue
	} else {
		limitPriceList = t.buyLimitQueue
		marketPriceList = t.buyMarketQueue
	}
	if exchangeOrder.Type == model.MarketPrice {
		//先处理市价 市价订单和限价的订单进行匹配
		t.matchMarketPriceWithLP(limitPriceList, exchangeOrder)
	} else { //限价单
		//限价单先和限价单进行成交，如果未成交或未完全成交，则继续与市价单进行成交
		t.matchLimitPriceWithLP(limitPriceList, exchangeOrder)
		if exchangeOrder.Status == model.Trading {
			t.matchLimitPriceWithMP(marketPriceList, exchangeOrder)
		}
		//证明还未交易完成，本轮交易订单未完全交易完成
		if exchangeOrder.Status == model.Trading {
			t.addLimitQueue(exchangeOrder)
			if exchangeOrder.Direction == model.BUY {
				t.sendTradPlateMsg(t.buyTradePlate)
			} else {
				t.sendTradPlateMsg(t.sellTradePlate)
			}
		}
	}
}

// 发送买卖盘的数据到kafka中
func (t *CoinTrade) sendTradPlateMsg(plate *TradePlate) {
	//发数据 数据格式是什么
	//这里result是前端需要展示的买卖盘数据格式,展示前24个格式
	result := plate.Result(24)
	marshal, _ := json.Marshal(result)
	data := database.KafkaData{
		Topic: "exchange_order_trade_plate",
		Key:   []byte(plate.Symbol),
		Data:  marshal,
	}
	//这里SendSync不一定能成功，可以写成重试三次，但是这里没有重试，因为放到kafka的是展示给前端的数据，失败了也没关系
	//因为交易量一般是比较大的，所以可以保证数据可以一定程度的实时刷新
	err := t.kafkaClient.SendSync(data)
	if err != nil {
		logx.Error(err)
	} else {
		logx.Info("======exchange_order_trade_plate send 成功....==========")
	}
}

func (t *CoinTrade) initData() {
	orderDomain := domain.NewExchangeOrderDomain(t.db)
	//应该去查询对应symbol的订单 将其赋值到coinTrade里面的各个队列中，同时加入买卖盘
	//这样可以在该symbol买卖盘创建的时候将之前的该symbol的订单都放入买卖盘中，不遗漏
	exchangeOrders, err := orderDomain.FindOrderListBySymbol(context.Background(), t.symbol, model.Trading)
	if err != nil {
		logx.Error(err)
		return
	}
	for _, v := range exchangeOrders {
		if v.Type == model.MarketPrice { //市价
			if v.Direction == model.BUY {
				t.bmMux.Lock()
				t.buyMarketQueue = append(t.buyMarketQueue, v)
				t.bmMux.Unlock()
				continue
			}
			if v.Direction == model.SELL {
				t.smMux.Lock()
				t.sellMarketQueue = append(t.sellMarketQueue, v)
				t.smMux.Unlock()
				continue
			}
			//市价单 不进入买卖盘的
		} else if v.Type == model.LimitPrice { //限价
			if v.Direction == model.BUY {
				t.buyLimitQueue.mux.Lock()
				//deal
				isPut := false
				for _, o := range t.buyLimitQueue.list {
					if o.price == v.Price {
						o.list = append(o.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					lpm := &LimitPriceMap{
						price: v.Price,
						list:  []*model.ExchangeOrder{v},
					}
					t.buyLimitQueue.list = append(t.buyLimitQueue.list, lpm)
				}
				t.buyTradePlate.Add(v) //放入买卖盘
				t.buyLimitQueue.mux.Unlock()
			} else if v.Direction == model.SELL {
				t.sellLimitQueue.mux.Lock()
				//deal
				isPut := false
				for _, o := range t.sellLimitQueue.list {
					if o.price == v.Price {
						o.list = append(o.list, v)
						isPut = true
						break
					}
				}
				if !isPut {
					lpm := &LimitPriceMap{
						price: v.Price,
						list:  []*model.ExchangeOrder{v},
					}
					t.sellLimitQueue.list = append(t.sellLimitQueue.list, lpm)
				}
				t.sellTradePlate.Add(v)
				t.sellLimitQueue.mux.Unlock()
			}
		}
	}
	//排序
	sort.Sort(t.buyMarketQueue)
	sort.Sort(t.sellMarketQueue)
	sort.Sort(t.buyLimitQueue.list)                //从高到低
	sort.Sort(sort.Reverse(t.sellLimitQueue.list)) //从低到高
	if len(exchangeOrders) > 0 {                   //有存量数据
		//将买盘和卖盘的数据推送到前端
		t.sendTradPlateMsg(t.buyTradePlate)
		t.sendTradPlateMsg(t.sellTradePlate)
	}
}

// matchLimitPriceWithMP 本限价单语市价队列进行匹配
// focusedOrder 限价单  mpList 市价队列
func (t *CoinTrade) matchLimitPriceWithMP(mpList TradeTimeQueue, focusedOrder *model.ExchangeOrder) {
	var delOrders []string
	for _, matchOrder := range mpList {
		if matchOrder.MemberId == focusedOrder.MemberId {
			//自己的订单就不处理了
			continue
		}
		price := focusedOrder.Price
		//计算可交易的数量
		matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
		if matchAmount <= 0 {
			continue
		}
		focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
		if matchAmount >= focusedAmount {
			//订单直接就交易完成了 能满足
			turnover := op.MulFloor(price, focusedAmount, 8)
			matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
			matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
			if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
				matchOrder.Status = model.Completed
				//从队列进行删除
				delOrders = append(delOrders, matchOrder.OrderId)
			}
			focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
			focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
			focusedOrder.Status = model.Completed
			break
		} else {
			//当前的订单 不满足交易额 继续进行匹配
			turnover := op.MulFloor(price, matchAmount, 8)
			matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
			matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
			matchOrder.Status = model.Completed
			//从队列进行删除
			delOrders = append(delOrders, matchOrder.OrderId)
			focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
			focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
			continue
		}
	}
	//处理已经匹配完成的订单 从队列删除
	for _, orderId := range delOrders {
		for index, matchOrder := range mpList {
			if matchOrder.OrderId == orderId {
				mpList = append(mpList[:index], mpList[index+1:]...)
				break
			}
		}
	}
}

// matchLimitPriceWithLP 限价单和限价队列中的订单进行匹配
// focusedOrder 限价单 lpList 限价队列
func (t *CoinTrade) matchLimitPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	var delOrders []string
	//标记买卖盘的数据是否更改
	buyNotify := false
	sellNotify := false
	var completeOrders []*model.ExchangeOrder
	//如果本订单是买，则匹配的卖队列的价格是从低到高
	//如果本订单是卖，则匹配的买队列的价格是从高到低
	//队列的排序在初始化的时候已经有序了
	for _, v := range lpList.list {
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				//自己的订单就不处理了（不自己买自己卖）
				continue
			}
			//如果本订单是买，因为卖队列的价格从低到高，所以如果买的价格比卖的价格还低，则无法成交
			if model.BUY == focusedOrder.Direction {
				if focusedOrder.Price < matchOrder.Price {
					break
				}
			}
			//如果是本订单是卖，因为买队列的价格从高到低，所以如果卖的价格比买的价格还高，则无法成交
			if model.SELL == focusedOrder.Direction {
				if focusedOrder.Price > matchOrder.Price {
					break
				}
			}
			//matchOrder和focusedOrder 是否匹配
			price := matchOrder.Price
			//计算可交易的数量
			matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}
			focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)
			if matchAmount >= focusedAmount {
				//订单直接就交易完成了 能满足
				turnover := op.MulFloor(price, focusedAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.Completed
					//从队列进行删除
					delOrders = append(delOrders, matchOrder.OrderId)
					completeOrders = append(completeOrders, matchOrder)
				}
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				focusedOrder.Status = model.Completed
				completeOrders = append(completeOrders, focusedOrder)
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focusedAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focusedAmount)
					sellNotify = true
				}
				break
			} else {
				//当前的订单 不满足交易额 继续进行匹配
				turnover := op.MulFloor(price, matchAmount, 8)
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				matchOrder.Status = model.Completed
				completeOrders = append(completeOrders, matchOrder)
				//从队列进行删除
				delOrders = append(delOrders, matchOrder.OrderId)

				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)

				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
	}
	//处理队列中 已经完成的订单进行删除
	for _, orderId := range delOrders {
		for _, v := range lpList.list {
			for index, matchOrder := range v.list {
				if orderId == matchOrder.OrderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	//通知买卖盘更新
	if buyNotify {
		t.sendTradPlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradPlateMsg(t.sellTradePlate)
	}
	for _, v := range completeOrders {
		t.sendCompleteOrder(v)
	}
}

// matchMarketPriceWithLP  市价订单交易函数
// focusedOrder 当前市价单  lpList 限价交易队列
func (t *CoinTrade) matchMarketPriceWithLP(lpList *LimitPriceQueue, focusedOrder *model.ExchangeOrder) {
	lpList.mux.Lock()
	defer lpList.mux.Unlock()
	//统计本次交易已完成的订单
	var delOrders []string
	//标记买卖盘的数据是否更改，如果更改了，则需要将买卖盘的最新数据发送给前端
	buyNotify := false
	sellNotify := false
	//如果是买，匹配卖的队列价格是从低到高，因为希望买的价格越低越好
	//如果是卖 买的队列价格是从高到低
	//队列的排序在初始化的时候已经排序好了（initData函数）
	for _, v := range lpList.list {
		for _, matchOrder := range v.list {
			if matchOrder.MemberId == focusedOrder.MemberId {
				//自己的订单就不处理了（不能自己买自己卖）
				continue
			}
			//matchOrder和focusedOrder 是否匹配
			price := matchOrder.Price
			//计算可交易的数量
			matchAmount := op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8)
			if matchAmount <= 0 {
				continue
			}
			focusedAmount := op.SubFloor(focusedOrder.Amount, focusedOrder.TradedAmount, 8)

			//市价单 买的时候amount的单位是USDT，这时候我们需要计算数量
			if focusedOrder.Direction == model.BUY {
				focusedAmount = op.DivFloor(op.SubFloor(focusedOrder.Amount, focusedOrder.Turnover, 8), price, 8)
			}
			if matchAmount >= focusedAmount { //能满足需求，订单直接就可以完全交易完成
				//计算成交额=交易的数量*当前队列订单的价格
				turnover := op.MulFloor(price, focusedAmount, 8)
				//修改队列中匹配到的订单的数据
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, focusedAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				//检查队列匹配到的订单是否已完成（这里只可能>=0，所以满足条件只可能Amount等于TradedAmount）
				if op.SubFloor(matchOrder.Amount, matchOrder.TradedAmount, 8) <= 0 {
					matchOrder.Status = model.Completed
					//从队列进行删除
					delOrders = append(delOrders, matchOrder.OrderId)
				}
				//修改本订单的数据
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, focusedAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)
				//修改本订单状态
				focusedOrder.Status = model.Completed
				//修改买卖盘中 匹配到的订单的相关数据
				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, focusedAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, focusedAmount)
					sellNotify = true
				}
				break
			} else {
				//当前的订单不满足交易额，需要继续进行匹配，所以这里后面用continue
				//计算成交额=交易的数量*当前队列订单的价格
				turnover := op.MulFloor(price, matchAmount, 8)
				//修改队列中匹配到的订单的数据和状态
				matchOrder.TradedAmount = op.AddFloor(matchOrder.TradedAmount, matchAmount, 8)
				matchOrder.Turnover = op.AddFloor(matchOrder.Turnover, turnover, 8)
				matchOrder.Status = model.Completed
				//从队列进行删除(因为此时队列中匹配到订单不够满足本订单的交易数量，所以本订单呗完全交易)
				delOrders = append(delOrders, matchOrder.OrderId)
				//修改本订单的数据
				focusedOrder.TradedAmount = op.AddFloor(focusedOrder.TradedAmount, matchAmount, 8)
				focusedOrder.Turnover = op.AddFloor(focusedOrder.Turnover, turnover, 8)

				if matchOrder.Direction == model.BUY {
					t.buyTradePlate.Remove(matchOrder, matchAmount)
					buyNotify = true
				} else {
					t.sellTradePlate.Remove(matchOrder, matchAmount)
					sellNotify = true
				}
				continue
			}
		}
	}
	//处理队列中 已经完成的订单进行删除
	for _, orderId := range delOrders {
		for _, v := range lpList.list {
			for index, matchOrder := range v.list {
				if orderId == matchOrder.OrderId {
					v.list = append(v.list[:index], v.list[index+1:]...)
					break
				}
			}
		}
	}
	//判断是否订单完成
	if focusedOrder.Status == model.Trading {
		//未完成 放入队列
		t.addMarketQueue(focusedOrder)
	}
	//websocket通知买卖盘更新
	if buyNotify {
		t.sendTradPlateMsg(t.buyTradePlate)
	}
	if sellNotify {
		t.sendTradPlateMsg(t.sellTradePlate)
	}
}

// 将未完成的市价订单加入市价订单队列中
func (t *CoinTrade) addMarketQueue(order *model.ExchangeOrder) {
	if order.Type != model.MarketPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyMarketQueue = append(t.buyMarketQueue, order)
		//处理完之后要重新排序
		sort.Sort(t.buyMarketQueue)
	} else {
		t.sellMarketQueue = append(t.sellMarketQueue, order)
		sort.Sort(t.sellMarketQueue)
	}
}

// 将限价订单添加到相应的队列中
func (t *CoinTrade) addLimitQueue(order *model.ExchangeOrder) {
	if order.Type != model.LimitPrice {
		return
	}
	if order.Direction == model.BUY {
		t.buyLimitQueue.mux.Lock()
		//deal
		isPut := false
		for _, o := range t.buyLimitQueue.list {
			if o.price == order.Price {
				o.list = append(o.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			lpm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.buyLimitQueue.list = append(t.buyLimitQueue.list, lpm)
		}
		t.buyTradePlate.Add(order)
		t.buyLimitQueue.mux.Unlock()
	} else if order.Direction == model.SELL {
		t.sellLimitQueue.mux.Lock()
		//deal
		isPut := false
		for _, o := range t.sellLimitQueue.list {
			if o.price == order.Price {
				o.list = append(o.list, order)
				isPut = true
				break
			}
		}
		if !isPut {
			lpm := &LimitPriceMap{
				price: order.Price,
				list:  []*model.ExchangeOrder{order},
			}
			t.sellLimitQueue.list = append(t.sellLimitQueue.list, lpm)
		}
		t.sellTradePlate.Add(order)
		t.sellLimitQueue.mux.Unlock()
	}
}

func (t *CoinTrade) sendCompleteOrder(order *model.ExchangeOrder) {
	if order.Status != model.Completed {
		return
	}
	marshal, _ := json.Marshal(order)
	kafkaData := database.KafkaData{
		Topic: "exchange_order_complete",
		Key:   []byte(t.symbol),
		Data:  marshal,
	}
	for {
		err := t.kafkaClient.SendSync(kafkaData)
		if err != nil {
			logx.Error(err)
			time.Sleep(250 * time.Millisecond)
			continue
		} else {
			break
		}
	}
}

func (t *CoinTrade) init() {
	t.buyTradePlate = NewTradePlate(t.symbol, model.BUY)
	t.sellTradePlate = NewTradePlate(t.symbol, model.SELL)
	t.buyLimitQueue = &LimitPriceQueue{}
	t.sellLimitQueue = &LimitPriceQueue{}
	t.initData()
}

func NewCoinTrade(symbol string, cli *database.KafkaClient, db *msdb.MsDB) *CoinTrade {
	c := &CoinTrade{
		symbol:      symbol,
		kafkaClient: cli,
		db:          db,
	}
	c.init()
	return c
}
