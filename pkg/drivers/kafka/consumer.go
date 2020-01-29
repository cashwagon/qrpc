package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/NightWolf007/qrpc/pkg/qrpc"
)

type kReader struct {
	*kafka.Reader
}

func (r *kReader) Init(cfg *kafka.ReaderConfig) {
	r.Reader = kafka.NewReader(*cfg)
}

// Consumer represents the wrapper on kafka.Reader.
// It implements the qrpc.Consumer interface.
type Consumer struct {
	TopicPrefix string
	cfg         *kafka.ReaderConfig
	r           Reader
}

// NewConsumer allocates new Consumer object
func NewConsumer(cfg *kafka.ReaderConfig) *Consumer {
	// Force disable auto-commit.
	// We use a manual commit for consistency.
	cfg.CommitInterval = 0

	return &Consumer{
		cfg: cfg,
		r:   &kReader{},
	}
}

// Subscribe subscribes consumer on multiple topics (queues) and starts it
// It must be called before Consume
func (c Consumer) Subscribe(queues []string) error {
	if len(queues) == 0 {
		return nil
	}

	// Package github.com/segmentio/kafka-go does not support multiple topics for now
	// So we pick only the first topic from the list
	// See: https://github.com/segmentio/kafka-go/issues/131
	c.cfg.Topic = queues[0]
	c.r.Init(c.cfg)

	return nil
}

// Consume runs one consume iteration.
// It fetches the message from the Kafka, calls qrpc.MessageHandler to process it
// and commits it
func (c Consumer) Consume(mh qrpc.MessageHandler) error {
	ctx := context.Background()

	msg, err := c.r.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("cannot read message: %w", err)
	}

	err = mh(qrpc.Message{
		Queue:  msg.Topic,
		Method: string(msg.Key),
		Data:   msg.Value,
	})
	if err != nil {
		return fmt.Errorf("cannot process message: %w", err)
	}

	if err := c.r.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("cannot commit message: %w", err)
	}

	return nil
}

// Close closes the consumer connection
func (c Consumer) Close() error {
	return c.r.Close()
}
