package qrpc

import (
	"context"
)

type Message struct {
	Queue  string
	Method string
	Data   []byte
}

type MessageHandler func(Message) error

type Producer interface {
	SetQueue(string)
	Produce(context.Context, Message) error
	Close() error
}

type Consumer interface {
	Subscribe([]string) error
	Consume(MessageHandler) error
	Close() error
}
