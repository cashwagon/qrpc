package qrpc

import (
	"fmt"
	"text/template"
)

const emptyType = "empty.Empty"

// templateData represents struct for text template
type templateData struct {
	Services   []Service
	ContextPkg string
	QrpcPkg    string
	UUIDPkg    string
}

// Service represent proto RPC service
type Service struct {
	Name                 string
	IsDeprecated         bool
	ForwardMethods       []Method
	BackwardMethods      []Method
	BidirectionalMethods []Method
}

// Method represents RPC method of proto service
type Method struct {
	Name         string
	IsDeprecated bool
	Comment      string
	InType       string
	OutType      string
}

// IsForward checks if this method is forward
func (m Method) IsForward() bool {
	return m.InType != emptyType && m.OutType == emptyType
}

// IsBackward checks if this method is backward
func (m Method) IsBackward() bool {
	return m.InType == emptyType && m.OutType != emptyType
}

// IsBidirectional checks if this method is bidirectional
func (m Method) IsBidirectional() bool {
	return m.InType != emptyType && m.OutType != emptyType
}

var templateFuncs = template.FuncMap{ // nolint:gochecknoglobals // there is no race
	"sprintf": fmt.Sprintf,
	"concat":  concat,
}

func concat(a, b []Method) []Method {
	return append(a, b...)
}
