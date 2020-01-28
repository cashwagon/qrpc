package test

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driver "github.com/NightWolf007/qrpc/pkg/drivers/kafka"
	"github.com/NightWolf007/qrpc/pkg/qrpc"
	"github.com/NightWolf007/qrpc/test/pb"
)

const (
	firstUid    = "12345"
	secondUid   = "54321"
	topicPrefix = "test"
	topic       = "test.qrpc.test.api.TestAPI"
)

type server struct {
	t    *testing.T
	done chan struct{}
}

func (s *server) FirstMethod(ctx context.Context, req *pb.FirstMethodRequest) error {
	assert.Equal(s.t, firstUid, req.Uid)
	s.done <- struct{}{}

	return nil
}

func (s *server) SecondMethod(ctx context.Context, req *pb.SecondMethodRequest) error {
	assert.Equal(s.t, secondUid, req.Uid)
	s.done <- struct{}{}

	return nil
}

func TestQRPCKafka(t *testing.T) {
	ctx := context.Background()
	brokers := []string{"kafka:9092"}

	// Create topic in kafka
	kconn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		strings.Join(brokers, ","),
		topic,
		0,
	)
	require.NoError(t, err, "cannot connect to kafka")
	defer kconn.Close()

	// Initialize server
	done := make(chan struct{})
	server := &server{
		t:    t,
		done: done,
	}

	// Intitialize and start server
	srv := qrpc.NewServer(
		driver.NewConsumer(&kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: "qrpc-test-group",
			Logger:  kafka.LoggerFunc(log.Printf),
		}),
		topicPrefix,
	)

	pb.RegisterTestAPIServer(srv, server)

	go func() {
		err := srv.Start()
		require.NoError(t, err, "unexpected server exit")
	}()
	defer srv.Stop()

	// Initialize client
	conn := qrpc.NewClientConn(
		driver.NewProducer(&kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
			Logger:   kafka.LoggerFunc(log.Printf),
		}),
		topicPrefix,
	)

	cli := pb.NewTestAPIClient(conn)

	// Send first request
	err = cli.FirstMethod(ctx, &pb.FirstMethodRequest{
		Uid: firstUid,
	})
	assert.NoError(t, err)

	// Send second request
	err = cli.SecondMethod(ctx, &pb.SecondMethodRequest{
		Uid: secondUid,
	})
	assert.NoError(t, err)

	// Close connection to flush producer buffer
	err = conn.Close()
	require.NoError(t, err, "cannot close producer connection")

	// Wait for processing two requests
	<-done
	<-done
}
