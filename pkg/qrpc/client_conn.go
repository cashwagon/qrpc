package qrpc

import (
	"context"
)

type ClientConn struct {
	qPrefix string
	p       Producer
}

func NewClientConn(p Producer, qPrefix string) *ClientConn {
	return &ClientConn{
		qPrefix: qPrefix,
		p:       p,
	}
}

func (c ClientConn) SetService(service string) {
	c.p.SetQueue(serviceToQueue(c.qPrefix, service))
}

func (c ClientConn) Invoke(ctx context.Context, msg Message) error {
	return c.p.Produce(ctx, msg)
}

func (c ClientConn) Close() error {
	return c.p.Close()
}
