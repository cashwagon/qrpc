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
)

const (
	groupID = "qrpc-examples-group"
)

type Server struct{}

func (s *Server) Echo(ctx context.Context, in *pb.EchoRequest) error {
	log.Printf("Received: %s", in.GetGreeting())
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

	pb.RegisterEchoAPIServer(srv, &Server{})

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}
	}()

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
