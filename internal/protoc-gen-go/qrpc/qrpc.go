// Package qrpc outputs qRPC service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-qrpcgo.
package qrpc

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// the qrpc package is introduced; the generated code references
// a constant, qrpc.SupportPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 1

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"
	qrpcPkgPath    = "github.com/NightWolf007/qrpc/pkg/qrpc"
)

// deprecationComment is the standard comment added to deprecated
// messages, fields, enums, and enum values.
const deprecationComment = "// Deprecated: Do not use."

func init() {
	generator.RegisterPlugin(new(qrpc))
}

// qrpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for qRPC support.
type qrpc struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "qrpc".
func (q *qrpc) Name() string {
	return "qrpc"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
	qrpcPkg    string
)

// Init initializes the plugin.
func (q *qrpc) Init(gen *generator.Generator) {
	q.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (q *qrpc) objectNamed(name string) generator.Object {
	q.gen.RecordTypeUse(name)
	return q.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (q *qrpc) typeName(str string) string {
	return q.gen.TypeName(q.objectNamed(str))
}

// P forwards to g.gen.P.
func (q *qrpc) P(args ...interface{}) {
	q.gen.P(args...)
}

// Generate generates code for the services in the given file.
func (q *qrpc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	contextPkg = string(q.gen.AddImport(contextPkgPath))
	qrpcPkg = string(q.gen.AddImport(qrpcPkgPath))

	// Assert version compatibility.
	q.P("// This is a compile-time assertion to ensure that this generated file")
	q.P("// is compatible with the qrpc package it is being compiled against.")
	q.P("const _ = ", qrpcPkg, ".SupportPackageIsVersion", generatedCodeVersion)
	q.P()

	for i, service := range file.FileDescriptorProto.Service {
		q.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (q *qrpc) GenerateImports(file *generator.FileDescriptor) {
}

// generateService generates all the code for the named service.
func (q *qrpc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServName := service.GetName()

	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}

	servName := generator.CamelCase(origServName)
	deprecated := service.GetOptions().GetDeprecated()

	q.P()
	q.P(fmt.Sprintf("// %sClient is the client API for %s service.", servName, servName))

	// Client interface.
	if deprecated {
		q.P("//")
		q.P(deprecationComment)
	}

	q.P("type ", servName, "Client interface {")

	for i, method := range service.Method {
		q.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.

		if method.GetOptions().GetDeprecated() {
			q.P("//")
			q.P(deprecationComment)
		}

		q.P(q.generateClientSignature(method))
	}

	q.P("}")
	q.P()

	// Client structure.
	q.P("type ", unexport(servName), "Client struct {")
	q.P("cc *", qrpcPkg, ".ClientConn")
	q.P("}")
	q.P()

	// NewClient factory.
	if deprecated {
		q.P(deprecationComment)
	}

	q.P("func New", servName, "Client (cc *", qrpcPkg, ".ClientConn) ", servName, "Client {")
	q.P("cc.SetService(", strconv.Quote(fullServName), ")")
	q.P("return &", unexport(servName), "Client{cc}")
	q.P("}")
	q.P()

	var methodIndex int

	serviceDescVar := "_" + servName + "_serviceDesc"

	// Client method implementations.
	for _, method := range service.Method {
		methodIndex++

		q.generateClientMethod(servName, method)
	}

	// Server interface.
	serverType := servName + "Server"
	q.P("// ", serverType, " is the server API for ", servName, " service.")

	if deprecated {
		q.P("//")
		q.P(deprecationComment)
	}

	q.P("type ", serverType, " interface {")

	for i, method := range service.Method {
		q.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.

		if method.GetOptions().GetDeprecated() {
			q.P("//")
			q.P(deprecationComment)
		}

		q.P(q.generateServerSignature(method))
	}

	q.P("}")
	q.P()

	// Server registration.
	if deprecated {
		q.P(deprecationComment)
	}

	q.P("func Register", servName, "Server(s *", qrpcPkg, ".Server, srv ", serverType, ") {")
	q.P("s.RegisterService(&", serviceDescVar, `, srv)`)
	q.P("}")
	q.P()

	// Server handler implementations.
	handlerNames := make([]string, 0, len(service.Method))

	for _, method := range service.Method {
		hname := q.generateServerMethod(servName, method)
		handlerNames = append(handlerNames, hname)
	}

	// Service descriptor.
	q.P("var ", serviceDescVar, " = ", qrpcPkg, ".ServiceDesc {")
	q.P("ServiceName: ", strconv.Quote(fullServName), ",")
	q.P("HandlerType: (*", serverType, ")(nil),")
	q.P("Methods: []", qrpcPkg, ".MethodDesc{")

	for i, method := range service.Method {
		q.P("{")
		q.P("MethodName: ", strconv.Quote(method.GetName()), ",")
		q.P("Handler: ", handlerNames[i], ",")
		q.P("},")
	}

	q.P("},")
	q.P("}")
	q.P()
}

// generateClientSignature returns the client-side signature for a method.
func (q *qrpc) generateClientSignature(method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)

	reqArg := ", in *" + q.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}

	return fmt.Sprintf("%s(ctx %s.Context%s) error", methName, contextPkg, reqArg)
}

func (q *qrpc) generateClientMethod(servName string, method *pb.MethodDescriptorProto) {
	if method.GetOptions().GetDeprecated() {
		q.P(deprecationComment)
	}

	q.P("func (c *", unexport(servName), "Client) ", q.generateClientSignature(method), "{")
	q.P("data, err := proto.Marshal(in)")
	q.P("if err != nil { return err }")
	q.P("return c.cc.Invoke(ctx, ", qrpcPkg, ".Message{")
	q.P("Method: ", strconv.Quote(method.GetName()), ",")
	q.P("Data: data,")
	q.P("})")
	q.P("}")
	q.P()
}

// generateServerSignature returns the server-side signature for a method.
func (q *qrpc) generateServerSignature(method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)

	reqArgs := []string{
		contextPkg + ".Context",
		"*" + q.typeName(method.GetInputType()),
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") error"
}

func (q *qrpc) generateServerMethod(servName string, method *pb.MethodDescriptorProto) string {
	methName := generator.CamelCase(method.GetName())
	hname := fmt.Sprintf("_%s_%s_Handler", servName, methName)
	inType := q.typeName(method.GetInputType())

	q.P("func ", hname, "(srv interface{}, ctx ", contextPkg, ".Context, msg []byte) error {")
	q.P("in := new(", inType, ")")
	q.P("if err := proto.Unmarshal(msg, in); err != nil { return err }")
	q.P("return srv.(", servName, "Server).", methName, "(ctx, in)")
	q.P("}")
	q.P()

	return hname
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
