package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/NightWolf007/qrpc/pkg/qrpc"
)

type kWriter struct {
	*kafka.Writer
}

func (w *kWriter) Init(cfg *kafka.WriterConfig) {
	w.Writer = kafka.NewWriter(*cfg)
}

type Producer struct {
	cfg *kafka.WriterConfig
	w   Writer
}

func NewProducer(cfg *kafka.WriterConfig) *Producer {
	return &Producer{
		cfg: cfg,
		w:   &kWriter{},
	}
}

func (p *Producer) SetQueue(queue string) {
	p.cfg.Topic = queue
	p.w.Init(p.cfg)
}

func (p Producer) Produce(ctx context.Context, msg qrpc.Message) error {
	err := p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Method),
		Value: msg.Data,
	})
	if err != nil {
		return fmt.Errorf("cannot write message: %w", err)
	}

	return nil
}

func (p Producer) Close() error {
	return p.w.Close()
}
