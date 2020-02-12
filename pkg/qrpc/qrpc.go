// Package qrpc provides qRPC implementation
package qrpc

import (
	"context"
	"fmt"
	"strings"
)

// The SupportPackageIsVersion variables are referenced from generated protocol
// buffer files to ensure compatibility with the qRPC version used.
// The latest support package version is 1.
// These constants should not be referenced from any other code.
const (
	SupportPackageIsVersion1 = true
)

type methodHandler func(srv interface{}, ctx context.Context, requestID string, msg []byte) error

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct {
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     []MethodDesc
}

type service struct {
	server interface{} // the server for service methods
	md     map[string]*MethodDesc
}

// serviceToQueue converts the service name to the queue name by joining it with prefix
func serviceToQueue(prefix, service string) string {
	if prefix == "" {
		return service
	}

	return strings.Join([]string{prefix, service}, ".")
}

// queueToService converts the queue name to the service name by removing prefix from it
func queueToService(prefix, queue string) string {
	if prefix == "" {
		return queue
	}

	return strings.Replace(queue, fmt.Sprintf("%s.", prefix), "", 1)
}
