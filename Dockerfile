FROM golang:1.23 as builder

WORKDIR /app

# Copy Go modules and tidy up dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o trading-chart-service ./cmd/server

# --- Runtime Stage ---
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/trading-chart-service .

USER nonroot:nonroot
ENTRYPOINT ["/trading-chart-service"]
    