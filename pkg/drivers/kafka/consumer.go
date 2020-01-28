package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/NightWolf007/qrpc/pkg/qrpc"
)

type Consumer struct {
	TopicPrefix string
	cfg         *kafka.ReaderConfig
	r           *kafka.Reader
}

func NewConsumer(cfg *kafka.ReaderConfig) *Consumer {
	cfg.CommitInterval = 0
	return &Consumer{cfg: cfg}
}

func (c *Consumer) Subscribe(queues []string) error {
	if len(queues) == 0 {
		return nil
	}

	// TODO: github.com/segmentio/kafka-go package does not support multiple topics for now
	// So we pick only the first topic from the list
	// See: https://github.com/segmentio/kafka-go/issues/131
	c.cfg.Topic = queues[0]
	c.r = kafka.NewReader(*c.cfg)

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
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("cannot close consumer: %w", err)
	}

	return nil
}
