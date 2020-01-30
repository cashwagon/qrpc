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
	defaultGreeting = "Hello, World!"
)

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
			log.Fatalf("cannot close connection: %v", err)
		}
	}()

	cli := pb.NewEchoAPIClient(conn)

	greeting := defaultGreeting
	if len(os.Args) > 1 {
		greeting = os.Args[1]
	}

	err := cli.Echo(context.Background(), &pb.EchoRequest{
		Greeting: greeting,
	})
	if err != nil {
		log.Fatalf("cannot send echo request: %v", err)
	}
}
