// qrpc package provides qRPC implementation
package qrpc

import (
	"context"
	"fmt"
	"strings"
)

const (
	SupportPackageIsVersion1 = true
)

type methodHandler func(srv interface{}, ctx context.Context, msg []byte) error

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

func serviceToQueue(prefix, service string) string {
	if prefix == "" {
		return service
	}

	return strings.Join([]string{prefix, service}, ".")
}

func queueToService(prefix, queue string) string {
	if prefix == "" {
		return queue
	}

	return strings.Replace(queue, fmt.Sprintf("%s.", prefix), "", 1)
}
