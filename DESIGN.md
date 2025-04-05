
---

## ✅ DESIGN.md (Architecture & Design)

### 📄 `DESIGN.md`

# Trading Chart Service – Design Document

## 🧱 Architecture Overview

### 🛰️ Data Ingestion
- Connects to Binance WebSocket API for `btcusdt`, `ethusdt`, `pepeusdt`
- Parses tick data and sends to the OHLC aggregator

### 📊 Aggregation
- Maintains 1-minute OHLC in-memory buffers per symbol
- Every minute, OHLC is flushed:
  - Broadcast to gRPC clients
  - Persisted to PostgreSQL

### 🧪 Streaming API
- gRPC `SubscribeCandles(symbol)` streams new candles in real time
- Supports multiple subscribers per symbol using a channel fan-out mechanism

### 💾 Persistence
- PostgreSQL schema optimized for time-series write performance
- Batched inserts per minute

### ☁️ Kubernetes & Infra
- Deployed to Kubernetes via Terraform
- Helm optional for PostgreSQL setup
- Can scale horizontally via more aggregator pods and external message queues

## 📈 Future Improvements
- Add Redis pub/sub for better fan-out
- Add Prometheus/Grafana monitoring
- Implement horizontal scaling using Kafka event bus
