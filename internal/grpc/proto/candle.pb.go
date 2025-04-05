package grpcserver

import (
	"log"
	"sync"

	pb "trading-chart-service/internal/grpc/proto"
	"trading-chart-service/internal/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

type CandleServer struct {
	pb.UnimplementedCandleServiceServer
	subscribers sync.Map // symbol -> list of client channels
}

func NewGRPCServer() *CandleServer {
	return &CandleServer{}
}

func (s *CandleServer) StartGRPCServer(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCandleServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s\n", address)
	return grpcServer.Serve(lis)
}

func (s *CandleServer) SubscribeCandles(req *pb.CandleRequest, stream pb.CandleService_SubscribeCandlesServer) error {
	symbol := req.Symbol
	ch := make(chan *pb.Candle, 10)

	// Add to subscribers list
	val, _ := s.subscribers.LoadOrStore(symbol, &[]chan *pb.Candle{})
	list := val.(*[]chan *pb.Candle)

	mu := sync.Mutex{}
	mu.Lock()
	*list = append(*list, ch)
	mu.Unlock()

	// Clean up on disconnect
	defer func() {
		mu.Lock()
		newList := make([]chan *pb.Candle, 0)
		for _, c := range *list {
			if c != ch {
				newList = append(newList, c)
			}
		}
		*list = newList
		mu.Unlock()
		close(ch)
		log.Printf("Client disconnected from %s", symbol)
	}()

	p, _ := peer.FromContext(stream.Context())
	log.Printf("Client %v subscribed to %s", p.Addr, symbol)

	for candle := range ch {
		if err := stream.Send(candle); err != nil {
			log.Printf("Send error: %v", err)
			break
		}
	}
	return nil
}

func (s *CandleServer) Broadcast(symbol string, ohlc models.OHLC) {
	val, ok := s.subscribers.Load(symbol)
	if !ok {
		return
	}

	candle := &pb.Candle{
		Symbol:    symbol,
		Open:      ohlc.Open,
		High:      ohlc.High,
		Low:       ohlc.Low,
		Close:     ohlc.Close,
		Volume:    ohlc.Volume,
		Timestamp: ohlc.Timestamp.Unix(),
	}

	for _, ch := range *val.(*[]chan *pb.Candle) {
		select {
		case ch <- candle:
		case <-time.After(100 * time.Millisecond):
			log.Printf("Slow client dropped")
		}
	}
}
