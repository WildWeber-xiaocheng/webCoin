package processor

import (
	"encoding/json"
	"market-api/internal/database"
	"market-api/internal/model"
)

const KLINE1M = "kline_1m"
const KLINE = "kline"
const TRADE = "trade"

type MarketHandler interface {
	HandlerTrade(symbol string, data []byte)        //处理交易数据
	HandlerKLine(symbol string, kline *model.Kline) //处理K线数据
}

type ProcessData struct {
	Type string //trade 交易 kline k线
	Key  []byte
	Data []byte
}

type Processor interface {
	Process(data ProcessData)
	AddHandler(h MarketHandler)
}

type DefaultProcessor struct {
	kafkaCli *database.KafkaClient
	handlers []MarketHandler
}

func NewDefaultProcessor(kafkaCli *database.KafkaClient) *DefaultProcessor {
	return &DefaultProcessor{
		kafkaCli: kafkaCli,
		handlers: make([]MarketHandler, 0),
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
			v.HandlerKLine(symbol, kline)
		}
	}
}

func (p *DefaultProcessor) Init() {
	//从kafka中接收kline 1m的同步数据
	p.startReadFromKafka(KLINE1M, KLINE)
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
