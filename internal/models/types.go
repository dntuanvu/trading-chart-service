package models

import (
	"strconv"
	"time"
)

type Trade struct {
	Price     float64
	Quantity  float64
	Timestamp time.Time
}

func ParseTrade(p, q string, t int64) (Trade, error) {
	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return Trade{}, err
	}

	quantity, err := strconv.ParseFloat(q, 64)
	if err != nil {
		return Trade{}, err
	}

	return Trade{
		Price:     price,
		Quantity:  quantity,
		Timestamp: time.UnixMilli(t),
	}, nil
}

type OHLC struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	Timestamp time.Time
}
