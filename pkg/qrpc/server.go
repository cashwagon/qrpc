package qrpc

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// Server represents the qRPC consumer
type Server struct {
	qPrefix string
	mu      sync.Mutex
	m       map[string]*service
	c       Consumer
	quit    chan struct{}
	done    chan struct{}
}

// NewServer allocates new Server struct
func NewServer(c Consumer, qPrefix string) *Server {
	return &Server{
		qPrefix: qPrefix,
		m:       make(map[string]*service),
		c:       c,
		quit:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// RegisterService registers a service and its implementation to the qRPC server.
// It is called from the IDL generated code. This must be called before invoking Start.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) error {
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(ss)

	if !st.Implements(ht) {
		return fmt.Errorf("found the handler of type %v that does not satisfy %v", st, ht)
	}

	return s.register(sd, ss)
}

// Start executes the consuming loop.
// In most cases it should be called in the separate goroutine.
func (s *Server) Start() error {
	queues := make([]string, 0, len(s.m))
	for sName := range s.m {
		queues = append(queues, serviceToQueue(s.qPrefix, sName))
	}

	if err := s.c.Subscribe(queues); err != nil {
		return fmt.Errorf("cannot subscribe on topics: %w", err)
	}

	for {
		select {
		case <-s.quit:
			s.done <- struct{}{}
			return nil
		default:
			if err := s.c.Consume(s.processMessage); err != nil {
				fmt.Printf("Consuming error: %v\n", err)
			}
		}
	}
}

// Stop sends the signal to the server to shutdown and waits until it's exit.
// The shutdown will happen only after processing the last consumed message.
func (s *Server) Stop() error {
	s.quit <- struct{}{}
	<-s.done

	if err := s.c.Close(); err != nil {
		return fmt.Errorf("cannot close server: %w", err)
	}

	return nil
}

func (s *Server) register(sd *ServiceDesc, ss interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[sd.ServiceName]; ok {
		return fmt.Errorf("found duplicate service registration for %q", sd.ServiceName)
	}

	srv := &service{
		server: ss,
		md:     make(map[string]*MethodDesc),
	}

	for i := range sd.Methods {
		d := &sd.Methods[i]
		srv.md[d.MethodName] = d
	}

	s.m[sd.ServiceName] = srv

	return nil
}

func (s *Server) processMessage(msg Message) error {
	service := queueToService(s.qPrefix, msg.Queue)

	srv, ok := s.m[service]
	if !ok {
		return fmt.Errorf("unknown service %s", msg.Queue)
	}

	md, ok := srv.md[msg.Method]
	if !ok {
		return fmt.Errorf("unknown method %s for service %s", msg.Method, msg.Queue)
	}

	return md.Handler(srv.server, context.Background(), msg.Data)
}
