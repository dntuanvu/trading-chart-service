package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"trading-chart-service/internal/models"
)

type StreamClient struct {
	Symbols   []string
	OnMessage func(symbol string, trade models.Trade)
}

func (c *StreamClient) Start(ctx context.Context) error {
	for _, symbol := range c.Symbols {
		go c.connectStream(ctx, symbol)
	}
	return nil
}

func (c *StreamClient) connectStream(ctx context.Context, symbol string) {
	u := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   fmt.Sprintf("/ws/%s@trade", symbol),
	}

	conn, _, err
