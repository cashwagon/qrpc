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

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) Init(cfg *kafka.ReaderConfig) {
	m.Called(cfg)
}

func (m *ReaderMock) FetchMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *ReaderMock) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	params := make([]interface{}, 0, len(msgs)+1)
	params = append(params, ctx)
	for _, msg := range msgs {
		params = append(params, msg)
	}

	args := m.Called(params...)
	return args.Error(0)
}

func (m *ReaderMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestConsumer_Subscribe(t *testing.T) {
	tests := []struct {
		name   string
		queues []string
		topic  string
	}{
		{
			"OneQueue",
			[]string{"queue1"},
			"queue1",
		},
		{
			"MultipleQueue",
			[]string{"queue1", "queue2", "queue3"},
			"queue1",
		},
		{
			"NoQueues",
			[]string{},
			"",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := &ReaderMock{}
			c := NewConsumer(&kafka.ReaderConfig{})
			c.r = r

			if tt.topic != "" {
				r.On("Init", &kafka.ReaderConfig{Topic: tt.topic}).Return().Once()
			}

			err := c.Subscribe(tt.queues)
			assert.NoError(t, err)

			r.AssertExpectations(t)
		})
	}
}

func TestConsumer_Consume(t *testing.T) {
	msg := kafka.Message{
		Topic:     "test.qrpc.Hello",
		Partition: 1,
		Key:       []byte("Hello"),
		Value:     []byte("testdata"),
	}

	t.Run("WhenSuccess", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("FetchMessage", mock.Anything).Return(msg, nil).Once()
		r.On("CommitMessages", mock.Anything, msg).Return(nil).Once()

		err := c.Consume(func(m qrpc.Message) error {
			assert.Equal(t, msg.Topic, m.Queue)
			assert.Equal(t, string(msg.Key), m.Method)
			assert.Equal(t, msg.Value, m.Data)
			return nil
		})
		assert.NoError(t, err)

		r.AssertExpectations(t)
	})

	t.Run("WhenFetchError", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("FetchMessage", mock.Anything).Return(msg, errors.New("fetch error")).Once()

		err := c.Consume(func(m qrpc.Message) error {
			assert.False(t, true, "MessageHandler was called")
			return nil
		})
		assert.Error(t, err)

		r.AssertExpectations(t)
	})

	t.Run("WhenMessageHanderError", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("FetchMessage", mock.Anything).Return(msg, nil).Once()

		err := c.Consume(func(m qrpc.Message) error {
			assert.Equal(t, msg.Topic, m.Queue)
			assert.Equal(t, string(msg.Key), m.Method)
			assert.Equal(t, msg.Value, m.Data)
			return errors.New("MessageHandler error")
		})
		assert.Error(t, err)

		r.AssertExpectations(t)
	})

	t.Run("WhenCommitError", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("FetchMessage", mock.Anything).Return(msg, nil).Once()
		r.On("CommitMessages", mock.Anything, msg).Return(errors.New("commit error")).Once()

		err := c.Consume(func(m qrpc.Message) error {
			assert.Equal(t, msg.Topic, m.Queue)
			assert.Equal(t, string(msg.Key), m.Method)
			assert.Equal(t, msg.Value, m.Data)
			return nil
		})
		assert.Error(t, err)

		r.AssertExpectations(t)
	})
}

func TestConsumer_Close(t *testing.T) {
	t.Run("WhenSuccess", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("Close").Return(nil).Once()

		err := c.Close()
		assert.NoError(t, err)

		r.AssertExpectations(t)
	})

	t.Run("WhenCloseFail", func(t *testing.T) {
		r := &ReaderMock{}
		c := NewConsumer(&kafka.ReaderConfig{})
		c.r = r

		r.On("Close").Return(errors.New("close error")).Once()

		err := c.Close()
		assert.Error(t, err)

		r.AssertExpectations(t)
	})
}
