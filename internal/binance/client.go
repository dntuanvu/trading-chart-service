package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/dntuanvu/trading-chart-service/internal/models"

	"github.com/gorilla/websocket"
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

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("WebSocket error for %s: %v", symbol, err)
		return
	}
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}

			var raw struct {
				P string `json:"p"` // Price
				Q string `json:"q"` // Quantity
				T int64  `json:"T"` // Trade time
			}

			if err := json.Unmarshal(msg, &raw); err != nil {
				log.Printf("Unmarshal error: %v", err)
				continue
			}

			trade, err := models.ParseTrade(raw.P, raw.Q, raw.T)
			if err != nil {
				log.Printf("Parse trade error: %v", err)
				continue
			}

			c.OnMessage(symbol, trade)
		}
	}
}
