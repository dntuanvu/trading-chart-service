# Trading Chart Service

A high-performance service that fetches real-time tick data from Binance, aggregates into OHLC candlesticks, streams via gRPC, and stores data in PostgreSQL.

## 🧩 Features
- Real-time Binance WebSocket integration
- 1-minute OHLC aggregation
- gRPC streaming API
- PostgreSQL persistence
- Docker + Kubernetes ready
- Terraform for infrastructure

## 🚀 Quick Start (Local Dev)

### Prerequisites
- Docker
- `protoc` (for gRPC codegen)

### Run locally

```bash
cd docker
docker-compose up --build
```

## 🧭 High-Level Architecture

```mermaid
flowchart TD
    Binance[Binance WebSocket API]
    Binance -->|Real-time Tick Data| GoService[Go Trading Service]

    GoService --> Aggregator[OHLC Aggregator (1-min)]
    GoService --> gRPC[gRPC Server]
    GoService --> Postgres[PostgreSQL DB]

    Aggregator -->|Every 1 min| gRPC
    gRPC --> Clients[Trader Dashboards (gRPC Subscribers)]

    subgraph Infrastructure
        K8s[Kubernetes Deployment]
        Terraform[Terraform IaC]
    end

    K8s --> GoService
    K8s --> Postgres