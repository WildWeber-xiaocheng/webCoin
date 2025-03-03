package processor

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/market/mclient"
	"grpc-common/market/types/market"
	"market-api/internal/database"
	"market-api/internal/model"
	"time"
)

const KLINE1M = "kline_1m"
const KLINE = "kline"
const TRADE = "trade"

type MarketHandler interface {
	HandlerTrade(symbol string, data []byte)                                               //处理交易数据
	HandlerKLine(symbol string, kline *model.Kline, thumbMap map[string]*market.CoinThumb) //处理K线数据
}

type ProcessData struct {
	Type string //trade 交易 kline k线
	Key  []byte
	Data []byte
}

type Processor interface {
	GetThumb() any
	Process(data ProcessData)
	AddHandler(h MarketHandler)
}

type DefaultProcessor struct {
	kafkaCli *database.KafkaClient
	handlers []MarketHandler
	thumbMap map[string]*market.CoinThumb
}

func NewDefaultProcessor(kafkaCli *database.KafkaClient) *DefaultProcessor {
	return &DefaultProcessor{
		kafkaCli: kafkaCli,
		handlers: make([]MarketHandler, 0),
		thumbMap: make(map[string]*market.CoinThumb),
	}
}

func (p *DefaultProcessor) AddHandler(h MarketHandler) {
	p.handlers = append(p.handlers, h)
}

func (p *DefaultProcessor) Process(data ProcessData) {
	if data.Type == KLINE {
		symbol := string(data.Key)
		kline := &model.Kline{}
		json.Unmarshal(data.Data, kline)
		for _, v := range p.handlers {
			v.HandlerKLine(symbol, kline, p.thumbMap)
		}
	}
}

func (p *DefaultProcessor) GetThumb() any {
	thumbs := make([]*market.CoinThumb, len(p.thumbMap))
	index := 0
	for _, thumb := range p.thumbMap {
		thumbs[index] = thumb
		index++
	}
	return thumbs
}

func (p *DefaultProcessor) Init(marketRpc mclient.Market) {
	//从kafka中接收kline 1m的同步数据
	p.startReadFromKafka(KLINE1M, KLINE)
	p.initThumbMap(marketRpc)
}

// startReadFromKafka 从kafka中接收kline 1m的同步数据
// typeTrade:处理数据的类型，值为kline/trade
func (p *DefaultProcessor) startReadFromKafka(topic string, typeTrade string) {
	//先start 后read
	p.kafkaCli.StartRead(topic)
	go p.dealQueueData(p.kafkaCli, typeTrade)
}

// dealQueueData 从kafka中真正读数据
func (p *DefaultProcessor) dealQueueData(cli *database.KafkaClient, typeTrade string) {
	for {
		kafkaData := cli.Read() //读kafka数据
		data := ProcessData{
			Type: typeTrade,
			Key:  kafkaData.Key,
			Data: kafkaData.Data,
		}
		p.Process(data)
	}
}

func (p *DefaultProcessor) initThumbMap(marketRpc mclient.Market) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	thumbResp, err := marketRpc.FindSymbolThumbTrend(ctx, &market.MarketReq{})
	if err != nil {
		logx.Info(err)
	} else {
		coinThumbs := thumbResp.List
		for _, thumb := range coinThumbs {
			p.thumbMap[thumb.Symbol] = thumb
		}
	}
}
