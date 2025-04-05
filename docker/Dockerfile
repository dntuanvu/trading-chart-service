FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o trading-chart-service ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/trading-chart-service .
COPY .env .

CMD ["./trading-chart-service"]