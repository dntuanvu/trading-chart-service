package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dntuanvu/trading-chart-service/internal/aggregator"
	"github.com/dntuanvu/trading-chart-service/internal/binance"
	grpcserver "github.com/dntuanvu/trading-chart-service/internal/grpc"

	"github.com/dntuanvu/trading-chart-service/internal/models"
	"github.com/dntuanvu/trading-chart-service/internal/persistence"
)

func main() {
	grpcSrv := grpcserver.NewGRPCServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := persistence.NewRepo(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	agg := aggregator.NewAggregator(func(symbol string, ohlc models.OHLC) {
		log.Printf("[FLUSH] %s => %+v\n", symbol, ohlc)

		// 1. Broadcast to gRPC
		grpcSrv.Broadcast(symbol, ohlc)

		// 2. Persist to DB (next step)
		repo.SaveOHLC(symbol, ohlc)
	})

	go grpcSrv.StartGRPCServer(":50051")

	client := &binance.StreamClient{
		Symbols: []string{"btcusdt", "ethusdt", "pepeusdt"},
		OnMessage: func(symbol string, trade models.Trade) {
			agg.AddTrade(symbol, trade)
		},
	}

	go client.Start(ctx)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			agg.FlushAll()
		}
	}()

	log.Println("Trading chart service started.")
	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down.")
}
