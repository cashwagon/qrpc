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
	pbh "github.com/cashwagon/qrpc/examples/pb/handler"
)

const (
	groupID = "qrpc-examples-handler-group"
)

type Server struct {
	cli pbh.EchoAPIClient
}

func (s *Server) Echo(ctx context.Context, reqID string, in *pb.EchoRequest) error {
	log.Printf("RequestID: %s - Received request: %s", reqID, in.GetGreeting())

	s.cli.Echo(ctx, reqID, &pb.EchoResponse{Greeting: in.GetGreeting()})
	return nil
}

func main() {
	brokersList := os.Getenv("KAFKA_BROKERS")
	if brokersList == "" {
		log.Fatal("Cannot get Kafka brokers list. Use KAFKA_BROKERS env variable")
	}

	brokers := strings.Split(brokersList, ",")

	conn := qrpc.NewClientConn(
		driver.NewProducer(&kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
		}),
		"examples",
	)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	cli := pbh.NewEchoAPIClient(conn)

	srv := qrpc.NewServer(
		driver.NewConsumer(&kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
		}),
		"examples",
	)

	pbh.RegisterEchoAPIServer(srv, &Server{cli})

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}
	}()

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
