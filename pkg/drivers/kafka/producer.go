package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/cashwagon/qrpc/pkg/qrpc"
)

type kWriter struct {
	*kafka.Writer
}

func (w *kWriter) Init(cfg *kafka.WriterConfig) {
	w.Writer = kafka.NewWriter(*cfg)
}

// Producer represents the wrapper on kafka.Writer.
// It implements the qrpc.Producer interface.
type Producer struct {
	cfg *kafka.WriterConfig
	w   Writer
}

// NewProducer allocates new Producer object
func NewProducer(cfg *kafka.WriterConfig) *Producer {
	return &Producer{
		cfg: cfg,
		w:   &kWriter{},
	}
}

// SetQueue sets the producer topic (queue) and starts the producer process.
// It must be called before Produce.
func (p Producer) SetQueue(queue string) {
	p.cfg.Topic = queue
	p.w.Init(p.cfg)
}

// Produce sends the message to the topic.
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

// Close closes the producer connection
func (p Producer) Close() error {
	return p.w.Close()
}
