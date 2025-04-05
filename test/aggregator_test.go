package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"trading-chart-service/internal/aggregator"
	"trading-chart-service/internal/models"
)

func TestAggregator_AddTrade(t *testing.T) {
	var flushed []models.OHLC

	agg := aggregator.NewAggregator(func(symbol string, ohlc models.OHLC) {
		flushed = append(flushed, ohlc)
	})

	ts := time.Now().Truncate(time.Minute)
	trades := []models.Trade{
		{Price: 100, Quantity: 1, Timestamp: ts},
		{Price: 105, Quantity: 2, Timestamp: ts.Add(10 * time.Second)},
		{Price: 99, Quantity: 1, Timestamp: ts.Add(20 * time.Second)},
		{Price: 103, Quantity: 3, Timestamp: ts.Add(50 * time.Second)},
	}

	for _, tr := range trades {
		agg.AddTrade("btcusdt", tr)
	}

	agg.FlushAll()

	assert.Len(t, flushed, 1)
	candle := flushed[0]
	assert.Equal(t, 100.0, candle.Open)
	assert.Equal(t, 105.0, candle.High)
	assert.Equal(t, 99.0, candle.Low)
	assert.Equal(t, 103.0, candle.Close)
	assert.Equal(t, 7.0, candle.Volume)
}
