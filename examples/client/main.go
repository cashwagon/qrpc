package main

import (
	"context"
	"log"
	"os"
	"strings"

	driver "github.com/cashwagon/qrpc/pkg/drivers/kafka"
	"github.com/cashwagon/qrpc/pkg/qrpc"
	"github.com/segmentio/kafka-go"

	"github.com/cashwagon/qrpc/examples/pb"
	pbc "github.com/cashwagon/qrpc/examples/pb/caller"
)

const (
	defaultGreeting = "Hello, World!"
	groupID         = "qrpc-examples-caller-group"
)

type Server struct {
	done chan struct{}
}

func NewServer() *Server {
	return &Server{make(chan struct{})}
}

func (s *Server) Echo(ctx context.Context, reqID string, out *pb.EchoResponse) error {
	log.Printf("RequestID: %s - Received response: %s", reqID, out.GetGreeting())
	s.done <- struct{}{}
	return nil
}

func main() {
	brokersList := os.Getenv("KAFKA_BROKERS")
	if brokersList == "" {
		log.Fatal("Cannot get Kafka brokers list. Use KAFKA_BROKERS env variable")
	}

	brokers := strings.Split(brokersList, ",")

	srv := qrpc.NewServer(
		driver.NewConsumer(&kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
		}),
		"examples",
	)

	s := NewServer()
	pbc.RegisterEchoAPIServer(srv, s)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}
	}()

	conn := qrpc.NewClientConn(
		driver.NewProducer(&kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
		}),
		"examples",
	)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Cannot close connection: %v", err)
		}
	}()

	cli := pbc.NewEchoAPIClient(conn)

	greeting := defaultGreeting
	if len(os.Args) > 1 {
		greeting = os.Args[1]
	}

	reqID, err := cli.Echo(context.Background(), &pb.EchoRequest{
		Greeting: greeting,
	})
	if err != nil {
		log.Fatalf("Cannot send echo request: %v", err)
	}

	log.Printf("RequestID: %s - Sent: %s", reqID, greeting)

	// Wait for response
	<-s.done
}
