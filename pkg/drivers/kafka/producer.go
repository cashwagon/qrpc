package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/NightWolf007/qrpc/pkg/qrpc"
)

type Producer struct {
	cfg *kafka.WriterConfig
	w   *kafka.Writer
}

func NewProducer(cfg *kafka.WriterConfig) *Producer {
	return &Producer{cfg: cfg}
}

func (p *Producer) SetQueue(queue string) {
	p.cfg.Topic = queue
	p.w = kafka.NewWriter(*p.cfg)
}

func (p Producer) Produce(ctx context.Context, msg qrpc.Message) error {
	err := p.w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.Method),
		Value: msg.Data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p Producer) Close() error {
	return p.w.Close()
}
