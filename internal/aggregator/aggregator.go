package aggregator

import (
	"sync"
	"time"

	"trading-chart-service/internal/models"
)

type Aggregator struct {
	mu      sync.Mutex
	ohlcMap map[string]*OHLCBuilder
	OnFlush func(symbol string, ohlc models.OHLC)
}

type OHLCBuilder struct {
	Symbol    string
	Current   models.OHLC
	startTime time.Time
}

func NewAggregator(flushFunc func(symbol string, ohlc models.OHLC)) *Aggregator {
	return &Aggregator{
		ohlcMap: make(map[string]*OHLCBuilder),
		OnFlush: flushFunc,
	}
}

func (a *Aggregator) AddTrade(symbol string, trade models.Trade) {
	a.mu.Lock()
	defer a.mu.Unlock()

	builder, exists := a.ohlcMap[symbol]
	if !exists || time.Since(builder.startTime) > time.Minute {
		if exists {
			a.OnFlush(symbol, builder.Current)
		}
		builder = &OHLCBuilder{
			Symbol: symbol,
			Current: models.OHLC{
				Open:      trade.Price,
				High:      trade.Price,
				Low:       trade.Price,
				Close:     trade.Price,
				Volume:    trade.Quantity,
				Timestamp: trade.Timestamp.Truncate(time.Minute),
			},
			startTime: trade.Timestamp.Truncate(time.Minute),
		}
		a.ohlcMap[symbol] = builder
		return
	}

	c := &builder.Current
	c.Close = trade.Price
	if trade.Price > c.High {
		c.High = trade.Price
	}
	if trade.Price < c.Low {
		c.Low = trade.Price
	}
	c.Volume += trade.Quantity
}

func (a *Aggregator) FlushAll() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for symbol, builder := range a.ohlcMap {
		a.OnFlush(symbol, builder.Current)
		delete(a.ohlcMap, symbol)
	}
}
