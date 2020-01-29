package qrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

type ProducerMock struct {
	mock.Mock
}

func (m *ProducerMock) SetQueue(queue string) {
	m.Called(queue)
}

func (m *ProducerMock) Produce(ctx context.Context, msg Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *ProducerMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestClientConn_SetService(t *testing.T) {
	p := &ProducerMock{}
	c := NewClientConn(p, "test")

	p.On("SetQueue", "test.qrpc.Hello").Return()

	c.SetService("qrpc.Hello")
	p.AssertExpectations(t)
}

func TestClientConn_Invoke(t *testing.T) {
	p := &ProducerMock{}
	c := NewClientConn(p, "test")

	ctx := context.Background()
	msg := Message{
		Method: "Hello",
		Data:   []byte("testdata"),
	}

	p.On("Produce", ctx, msg).Return(nil)

	c.Invoke(ctx, msg)
	p.AssertExpectations(t)
}

func TestClientConn_Close(t *testing.T) {
	p := &ProducerMock{}
	c := NewClientConn(p, "test")

	p.On("Close").Return(nil)

	c.Close()
	p.AssertExpectations(t)
}
