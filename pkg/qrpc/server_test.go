package qrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ConsumerMock struct {
	mock.Mock
}

func (m *ConsumerMock) Subscribe(queues []string) error {
	args := m.Called(queues)
	return args.Error(0)
}

func (m *ConsumerMock) Consume(mh MessageHandler) error {
	args := m.Called(mh)
	return args.Error(0)
}

func (m *ConsumerMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

type HelloRequest struct{}

type HelloServer interface {
	Hello(context.Context, *HelloRequest) error
}

type HelloServerMock struct {
	mock.Mock
}

func (m *HelloServerMock) Hello(ctx context.Context, req *HelloRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestServer_RegisterService(t *testing.T) {
	s := NewServer(nil, "")

	type args struct {
		sd *ServiceDesc
		ss interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Success",
			args{
				sd: &ServiceDesc{
					ServiceName: "qrpc.Hello",
					HandlerType: (*HelloServer)(nil),
					Methods: []MethodDesc{
						{
							MethodName: "Hello1",
							Handler:    nil,
						},
						{
							MethodName: "Hello2",
							Handler:    nil,
						},
					},
				},
				ss: &HelloServerMock{},
			},
			false,
		},
		{
			"WhenServerNotImplementHandler",
			args{
				sd: &ServiceDesc{
					ServiceName: "qrpc.Hello",
					HandlerType: (*HelloServer)(nil),
					Methods: []MethodDesc{
						{
							MethodName: "Hello1",
							Handler:    nil,
						},
						{
							MethodName: "Hello2",
							Handler:    nil,
						},
					},
				},
				ss: &struct{}{},
			},
			true,
		},
		{
			"WhenDuplicateService",
			args{
				sd: &ServiceDesc{
					ServiceName: "qrpc.Hello",
					HandlerType: (*HelloServer)(nil),
					Methods: []MethodDesc{
						{
							MethodName: "Hello1",
							Handler:    nil,
						},
					},
				},
				ss: &HelloServerMock{},
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := s.RegisterService(tt.args.sd, tt.args.ss)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !tt.wantErr {
				srv, ok := s.m[tt.args.sd.ServiceName]
				assert.True(t, ok)
				assert.Equal(t, tt.args.ss, srv.server)

				for _, md := range tt.args.sd.Methods {
					m, ok := srv.md[md.MethodName]
					assert.True(t, ok)
					assert.Equal(t, md, *m)
				}
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("WhenSuccess", func(t *testing.T) {
		c := &ConsumerMock{}
		msg := Message{
			Queue:  "test.qrpc.Hello",
			Method: "Hello",
			Data:   []byte("testdata"),
		}
		s := NewServer(c, "test")

		err := s.RegisterService(
			&ServiceDesc{
				ServiceName: "qrpc.Hello",
				HandlerType: (*HelloServer)(nil),
				Methods: []MethodDesc{
					{
						MethodName: "Hello",
						Handler: func(srv interface{}, ctx context.Context, m []byte) error {
							assert.Equal(t, msg.Data, m)
							return nil
						},
					},
				},
			},
			&HelloServerMock{},
		)
		require.NoError(t, err)

		ready := make(chan struct{})

		c.On("Subscribe", []string{"test.qrpc.Hello"}).Return(nil).Once()

		c.On("Consume", mock.AnythingOfType("MessageHandler")).Run(func(args mock.Arguments) {
			assert.NoError(t, args.Get(0).(MessageHandler)(msg))
			ready <- struct{}{}
		}).Return(nil).Once()
		c.On("Consume", mock.AnythingOfType("MessageHandler")).Return(nil)

		c.On("Close").Return(nil)

		go func() {
			assert.NoError(t, s.Start())
		}()

		<-ready

		err = s.Stop()
		assert.NoError(t, err)

		c.AssertExpectations(t)
	})

	t.Run("WhenSubscribeFail", func(t *testing.T) {
		c := &ConsumerMock{}
		s := NewServer(c, "test")

		err := s.RegisterService(
			&ServiceDesc{
				ServiceName: "qrpc.Hello",
				HandlerType: (*HelloServer)(nil),
				Methods: []MethodDesc{
					{
						MethodName: "Hello",
						Handler:    nil,
					},
				},
			},
			&HelloServerMock{},
		)
		require.NoError(t, err)

		c.On("Subscribe", []string{"test.qrpc.Hello"}).Return(errors.New("subscribe error")).Once()

		err = s.Start()
		assert.Error(t, err)
		c.AssertExpectations(t)
	})

	t.Run("WhenConsumeFail", func(t *testing.T) {
		c := &ConsumerMock{}
		s := NewServer(c, "test")

		err := s.RegisterService(
			&ServiceDesc{
				ServiceName: "qrpc.Hello",
				HandlerType: (*HelloServer)(nil),
				Methods: []MethodDesc{
					{
						MethodName: "Hello",
						Handler:    nil,
					},
				},
			},
			&HelloServerMock{},
		)
		require.NoError(t, err)

		ready := make(chan struct{})

		c.On("Subscribe", []string{"test.qrpc.Hello"}).Return(nil).Once()

		c.On("Consume", mock.AnythingOfType("MessageHandler")).Run(func(args mock.Arguments) {
			ready <- struct{}{}
		}).Return(errors.New("consume error")).Once()
		c.On("Consume", mock.AnythingOfType("MessageHandler")).Return(nil)

		c.On("Close").Return(nil)

		go func() {
			assert.NoError(t, s.Start())
		}()

		<-ready

		err = s.Stop()
		require.NoError(t, err)

		c.AssertExpectations(t)
	})
}

func TestServer_processMessage(t *testing.T) {
	msg := Message{
		Queue:  "test.qrpc.Hello",
		Method: "Hello",
		Data:   []byte("testdata"),
	}

	tests := []struct {
		name    string
		sd      ServiceDesc
		wantErr bool
	}{
		{
			"HandlerFound",
			ServiceDesc{
				ServiceName: "qrpc.Hello",
				HandlerType: (*HelloServer)(nil),
				Methods: []MethodDesc{
					{
						MethodName: "Hello",
						Handler: func(srv interface{}, ctx context.Context, m []byte) error {
							assert.Equal(t, msg.Data, m)
							return nil
						},
					},
				},
			},
			false,
		},
		{
			"MethodNotFound",
			ServiceDesc{
				ServiceName: "qrpc.Hello",
				HandlerType: (*HelloServer)(nil),
				Methods: []MethodDesc{
					{
						MethodName: "NotHello",
						Handler: func(srv interface{}, ctx context.Context, m []byte) error {
							assert.Equal(t, msg.Data, m)
							return nil
						},
					},
				},
			},
			true,
		},
		{
			"ServiceNotFound",
			ServiceDesc{
				ServiceName: "qrpc.NotHello",
				HandlerType: (*HelloServer)(nil),
				Methods:     []MethodDesc{},
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(nil, "test")
			err := s.RegisterService(&tt.sd, &HelloServerMock{})
			require.NoError(t, err)

			err = s.processMessage(msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
