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
    GoService[Go Trading Chart Service]
    Aggregator[OHLC Aggregator (1-min)]
    GRPC[gRPC Server]
    Postgres[PostgreSQL DB]
    Clients[gRPC Clients (Dashboards, Subscribers)]

    Binance --> GoService
    GoService --> Aggregator
    GoService --> GRPC
    GoService --> Postgres
    Aggregator --> GRPC
    GRPC --> Clients

    subgraph Infrastructure
        K8s[Kubernetes Cluster]
        Terraform[Terraform IaC]
    end

    K8s --> GoService
    K8s --> Postgres
