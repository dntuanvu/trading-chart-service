package persistence

import (
	"context"
	"log"

	"github.com/dntuanvu/trading-chart-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(dsn string) (*Repo, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

func (r *Repo) SaveOHLC(symbol string, ohlc models.OHLC) {
	_, err := r.db.Exec(context.Background(), `
		INSERT INTO candles (symbol, timestamp, open, high, low, close, volume)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, symbol, ohlc.Timestamp, ohlc.Open, ohlc.High, ohlc.Low, ohlc.Close, ohlc.Volume)
	if err != nil {
		log.Printf("DB insert error: %v", err)
	}
}

func (r *Repo) Close() {
	r.db.Close()
}
