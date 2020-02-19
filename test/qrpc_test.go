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
)

const (
	requestUID  = "12345"
	responseUID = "54321"
	topicPrefix = "test"
)

var (
	forwardReqID       string
	backwardReqID      string
	bidirectionalReqID string
)

type handlerServer struct {
	t       *testing.T
	brokers []string
	done    chan struct{}
}

func (s *handlerServer) BidirectionalMethod(ctx context.Context, reqID string, req *pb.Request) error {
	assert.Equal(s.t, bidirectionalReqID, reqID)
	assert.Equal(s.t, requestUID, req.GetUid())

	// Initialize handler client
	conn := initQRPCClientConn(s.brokers)

	defer func() {
		require.NoError(s.t, conn.Close(), "cannot close handler client connection")
	}()

	cli := pb.NewHandlerTestAPIClient(conn)

	// Send response
	err := cli.BidirectionalMethod(ctx, reqID, &pb.Response{Uid: responseUID})
	require.NoError(s.t, err)

	// Call backward method
	backwardReqID, err = cli.BackwardMethod(ctx, &pb.Response{Uid: responseUID})
	require.NoError(s.t, err)

	s.done <- struct{}{}

	return nil
}

func (s *handlerServer) ForwardMethod(ctx context.Context, reqID string, req *pb.Request) error {
	assert.Equal(s.t, forwardReqID, reqID)
	assert.Equal(s.t, requestUID, req.GetUid())

	s.done <- struct{}{}

	return nil
}

type callerServer struct {
	t    *testing.T
	done chan struct{}
}

func (s *callerServer) BidirectionalMethod(ctx context.Context, reqID string, resp *pb.Response) error {
	assert.Equal(s.t, bidirectionalReqID, reqID)
	assert.Equal(s.t, responseUID, resp.GetUid())

	s.done <- struct{}{}

	return nil
}

func (s *callerServer) BackwardMethod(ctx context.Context, reqID string, resp *pb.Response) error {
	assert.Equal(s.t, backwardReqID, reqID)
	assert.Equal(s.t, responseUID, resp.GetUid())

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
	pb.RegisterHandlerTestAPIServer(hsrv, hs)

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
	pb.RegisterCallerTestAPIServer(csrv, cs)

	go func() {
		require.NoError(t, csrv.Start(), "unexpected caller server exit")
	}()

	defer func() {
		require.NoError(t, csrv.Stop(), "cannot stop caller server")
	}()

	// Initialize caller client
	conn := initQRPCClientConn(brokers)
	cli := pb.NewCallerTestAPIClient(conn)

	var err error

	// Send first request
	forwardReqID, err = cli.ForwardMethod(ctx, &pb.Request{Uid: requestUID})
	assert.NoError(t, err)
	assert.NotEmpty(t, forwardReqID)

	// Send second request
	bidirectionalReqID, err = cli.BidirectionalMethod(ctx, &pb.Request{Uid: requestUID})
	assert.NoError(t, err)
	assert.NotEmpty(t, bidirectionalReqID)

	// Close connection to flush producer buffer
	err = conn.Close()
	require.NoError(t, err, "cannot close caller client connection")

	// Wait for processing requests
	<-hs.done
	<-hs.done
	<-cs.done
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
