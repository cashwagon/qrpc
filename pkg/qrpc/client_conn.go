package qrpc

import (
	"context"
)

// ClientConn represents producer connection to the message broker
type ClientConn struct {
	qPrefix string
	p       Producer
}

// NewClientConn creates new ClientConn object
func NewClientConn(p Producer, qPrefix string) *ClientConn {
	return &ClientConn{
		qPrefix: qPrefix,
		p:       p,
	}
}

// SetService sets the message broker queue based on the given service name
func (c ClientConn) SetService(service string) {
	c.p.SetQueue(serviceToQueue(c.qPrefix, service))
}

// Invoke sends the message to the queue in message broker
func (c ClientConn) Invoke(ctx context.Context, msg Message) error {
	return c.p.Produce(ctx, msg)
}

// Close closes the producer connection
func (c ClientConn) Close() error {
	return c.p.Close()
}
