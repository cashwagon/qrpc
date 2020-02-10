package test

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driver "github.com/cashwagon/qrpc/pkg/drivers/kafka"
	"github.com/cashwagon/qrpc/pkg/qrpc"
	"github.com/cashwagon/qrpc/test/pb"
	"github.com/cashwagon/qrpc/test/pb/caller"
	"github.com/cashwagon/qrpc/test/pb/handler"
)

const (
	firstUID    = "12345"
	secondUID   = "54321"
	topicPrefix = "test"
)

type handlerServer struct {
	t       *testing.T
	brokers []string
	done    chan struct{}
}

func (s *handlerServer) UnaryMethod(ctx context.Context, req *pb.UnaryMethodRequest) error {
	assert.Equal(s.t, firstUID, req.GetUid())
	s.done <- struct{}{}

	return nil
}

func (s *handlerServer) BinaryMethod(ctx context.Context, req *pb.BinaryMethodRequest) error {
	assert.Equal(s.t, secondUID, req.GetUid())

	// Initialize handler client
	conn := initQRPCClientConn(s.brokers)

	defer func() {
		require.NoError(s.t, conn.Close(), "cannot close handler client connection")
	}()

	cli := handler.NewTestAPIClient(conn)

	// Send response
	err := cli.BinaryMethod(ctx, &pb.BinaryMethodResponse{
		Uid: req.GetUid(),
	})
	require.NoError(s.t, err)

	s.done <- struct{}{}

	return nil
}

type callerServer struct {
	t    *testing.T
	done chan struct{}
}

func (s *callerServer) BinaryMethod(ctx context.Context, resp *pb.BinaryMethodResponse) error {
	assert.Equal(s.t, secondUID, resp.GetUid())
	s.done <- struct{}{}

	return nil
}

func TestQRPCKafka(t *testing.T) {
	ctx := context.Background()

	brokersList := os.Getenv("KAFKA_BROKERS")
	require.NotEmpty(t, brokersList, "Cannot get Kafka brokers list. Use KAFKA_BROKERS env variable")

	brokers := strings.Split(brokersList, ",")

	// Create topics in kafka
	kafkaCreateTopic(t, brokers, "test.qrpc.test.api.TestAPI.in")
	kafkaCreateTopic(t, brokers, "test.qrpc.test.api.TestAPI.out")

	// Initialize handler server
	hs := &handlerServer{
		t:       t,
		brokers: brokers,
		done:    make(chan struct{}),
	}

	// Intitialize and start handler server
	hsrv := initQRPCServer(brokers, "qrpc-test-group-handler")
	handler.RegisterTestAPIServer(hsrv, hs)

	go func() {
		require.NoError(t, hsrv.Start(), "unexpected handler server exit")
	}()

	defer func() {
		require.NoError(t, hsrv.Stop(), "cannot stop handler server")
	}()

	// Initialize caller server
	cs := &callerServer{
		t:    t,
		done: make(chan struct{}),
	}

	// Initialize and start caller server
	csrv := initQRPCServer(brokers, "qrpc-test-group-caller")
	caller.RegisterTestAPIServer(csrv, cs)

	go func() {
		require.NoError(t, csrv.Start(), "unexpected caller server exit")
	}()

	defer func() {
		require.NoError(t, csrv.Stop(), "cannot stop caller server")
	}()

	// Initialize caller client
	conn := initQRPCClientConn(brokers)
	cli := caller.NewTestAPIClient(conn)

	// Send first request
	err := cli.UnaryMethod(ctx, &pb.UnaryMethodRequest{
		Uid: firstUID,
	})
	assert.NoError(t, err)

	// Send second request
	err = cli.BinaryMethod(ctx, &pb.BinaryMethodRequest{
		Uid: secondUID,
	})
	assert.NoError(t, err)

	// Close connection to flush producer buffer
	err = conn.Close()
	require.NoError(t, err, "cannot close caller client connection")

	// Wait for processing requests
	<-hs.done
	<-hs.done
	<-cs.done
}

func kafkaCreateTopic(t *testing.T, brokers []string, topic string) {
	t.Helper()

	kconn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		strings.Join(brokers, ","),
		topic,
		0,
	)
	require.NoError(t, err, "cannot connect to kafka")
	require.NoError(t, kconn.Close())
}

func initQRPCServer(brokers []string, group string) *qrpc.Server {
	return qrpc.NewServer(
		driver.NewConsumer(&kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: group,
			Logger:  kafka.LoggerFunc(log.Printf),
		}),
		topicPrefix,
	)
}

func initQRPCClientConn(brokers []string) *qrpc.ClientConn {
	return qrpc.NewClientConn(
		driver.NewProducer(&kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
			Logger:   kafka.LoggerFunc(log.Printf),
		}),
		topicPrefix,
	)
}
