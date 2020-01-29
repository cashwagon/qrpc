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

type Consumer struct {
	TopicPrefix string
	cfg         *kafka.ReaderConfig
	r           Reader
}

func NewConsumer(cfg *kafka.ReaderConfig) *Consumer {
	cfg.CommitInterval = 0
	return &Consumer{
		cfg: cfg,
		r:   &kReader{},
	}
}

func (c *Consumer) Subscribe(queues []string) error {
	if len(queues) == 0 {
		return nil
	}

	// TODO: github.com/segmentio/kafka-go package does not support multiple topics for now
	// So we pick only the first topic from the list
	// See: https://github.com/segmentio/kafka-go/issues/131
	c.cfg.Topic = queues[0]
	c.r.Init(c.cfg)

	return nil
}

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

func (c Consumer) Close() error {
	return c.r.Close()
}
