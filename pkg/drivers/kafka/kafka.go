// kafka package implements Apache Kafka driver for qRPC
package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Reader represents abstract reader interface to wrap kafka.Reader
type Reader interface {
	Init(cfg *kafka.ReaderConfig)
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Writer represents abstract writer interface to wrap kafka.Writer
type Writer interface {
	Init(cfg *kafka.WriterConfig)
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}
