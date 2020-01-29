package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/NightWolf007/qrpc/pkg/qrpc"
)

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) Init(cfg *kafka.WriterConfig) {
	m.Called(cfg)
}

func (m *WriterMock) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	params := make([]interface{}, 0, len(msgs)+1)
	params = append(params, ctx)

	for i := range msgs {
		params = append(params, msgs[i])
	}

	args := m.Called(params...)

	return args.Error(0)
}

func (m *WriterMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestProducer_SetQueue(t *testing.T) {
	w := &WriterMock{}
	p := NewProducer(&kafka.WriterConfig{})
	p.w = w

	w.On("Init", &kafka.WriterConfig{Topic: "test.qrpc.Hello"}).Return().Once()

	p.SetQueue("test.qrpc.Hello")
	w.AssertExpectations(t)
}

func TestProducer_Produce(t *testing.T) {
	ctx := context.Background()
	msg := qrpc.Message{
		Method: "Hello",
		Data:   []byte("testdata"),
	}

	t.Run("WhenSuccess", func(t *testing.T) {
		w := &WriterMock{}
		p := NewProducer(&kafka.WriterConfig{})
		p.w = w

		w.On("WriteMessages", ctx, kafka.Message{
			Key:   []byte(msg.Method),
			Value: msg.Data,
		}).Return(nil).Once()

		err := p.Produce(ctx, msg)
		assert.NoError(t, err)

		w.AssertExpectations(t)
	})

	t.Run("WhenWriteError", func(t *testing.T) {
		w := &WriterMock{}
		p := NewProducer(&kafka.WriterConfig{})
		p.w = w

		w.On("WriteMessages", ctx, kafka.Message{
			Key:   []byte(msg.Method),
			Value: msg.Data,
		}).Return(errors.New("write error")).Once()

		err := p.Produce(ctx, msg)
		assert.Error(t, err)

		w.AssertExpectations(t)
	})
}

func TestProducer_Close(t *testing.T) {
	t.Run("WhenSuccess", func(t *testing.T) {
		w := &WriterMock{}
		p := NewProducer(&kafka.WriterConfig{})
		p.w = w

		w.On("Close").Return(nil).Once()

		err := p.Close()
		assert.NoError(t, err)

		w.AssertExpectations(t)
	})

	t.Run("WhenCloseFail", func(t *testing.T) {
		w := &WriterMock{}
		p := NewProducer(&kafka.WriterConfig{})
		p.w = w

		w.On("Close").Return(errors.New("close error")).Once()

		err := p.Close()
		assert.Error(t, err)

		w.AssertExpectations(t)
	})
}
