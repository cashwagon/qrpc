package qrpc

import (
	"context"
)

// Message represents the qRPC message struct
type Message struct {
	Queue     string
	Method    string
	RequestID string
	Data      []byte
}

// MessageHandler is a func, that process the message received from message broker
type MessageHandler func(Message) error

// Producer represents an abstract producer
// It should be implemented by the driver
type Producer interface {
	SetQueue(queue string)
	Produce(ctx context.Context, msg Message) error
	Close() error
}

// Consumer represents an abstract consumer
// It should be implemented by the driver
type Consumer interface {
	Subscribe(queues []string) error
	Consume(mh MessageHandler) error
	Close() error
}
