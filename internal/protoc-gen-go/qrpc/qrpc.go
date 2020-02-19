// Package qrpc outputs qRPC service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.
package qrpc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

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
	qrpcPkgPath    = "github.com/cashwagon/qrpc/pkg/qrpc"
	emptyPkgPath   = "github.com/golang/protobuf/ptypes/empty"
	uuidPkgPath    = "github.com/google/uuid"
)

func init() {
	generator.RegisterPlugin(new(qrpc))
}

// qrpc is an implementation of the Go protocol buffer compiler's
// plugin architecture. It generates bindings for qRPC support.
type qrpc struct {
	gen         *generator.Generator
	callerTmpl  *template.Template
	handlerTmpl *template.Template
	comments    map[string]*pb.SourceCodeInfo_Location
}

// Name returns the name of this plugin, "qrpc".
func (q *qrpc) Name() string {
	return "qrpc"
}

// Init initializes the plugin.
func (q *qrpc) Init(gen *generator.Generator) {
	q.gen = gen

	var err error

	q.callerTmpl, err = template.New("caller").Funcs(templateFuncs).Parse(callerTemplate)
	if err != nil {
		q.gen.Error(err, "failed to inititalize caller template")
	}

	q.handlerTmpl, err = template.New("handler").Funcs(templateFuncs).Parse(handlerTemplate)
	if err != nil {
		q.gen.Error(err, "failed to inititalize handler template")
	}
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
func (q *qrpc) P(args ...interface{}) { q.gen.P(args...) }

// Generate generates code for the services in the given file.
func (q *qrpc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	if err := q.generate(file.FileDescriptorProto); err != nil {
		q.gen.Error(err, "failed to generate file")
	}
}

// GenerateImports generates the import declaration for this file.
func (q *qrpc) GenerateImports(file *generator.FileDescriptor) {
}

func (q *qrpc) generate(file *pb.FileDescriptorProto) error {
	q.extractComments(file)

	contextPkg := string(q.gen.AddImport(contextPkgPath))
	qrpcPkg := string(q.gen.AddImport(qrpcPkgPath))
	emptyPkg := string(q.gen.AddImport(emptyPkgPath))
	uuidPkg := string(q.gen.AddImport(uuidPkgPath))

	data := templateData{
		Services:   q.generateServices(file),
		ContextPkg: contextPkg,
		QrpcPkg:    qrpcPkg,
		UUIDPkg:    uuidPkg,
	}

	b := new(bytes.Buffer)
	if err := q.renderTemplates(b, data); err != nil {
		return err
	}

	q.P("// Reference imports to suppress errors if they are not otherwise used.")
	q.P("var _ ", contextPkg, ".Context")
	q.P("var _ ", qrpcPkg, ".ClientConn")
	q.P("var _ ", emptyPkg, ".Empty")
	q.P("var _ ", uuidPkg, ".UUID")
	q.P()

	// Assert version compatibility.
	q.P("// This is a compile-time assertion to ensure that this generated file")
	q.P("// is compatible with the qrpc package it is being compiled against.")
	q.P("const _ = ", "qrpc.SupportPackageIsVersion", generatedCodeVersion)
	q.P()

	q.P(b.String())

	return nil
}

func (q *qrpc) generateServices(file *pb.FileDescriptorProto) []Service {
	services := make([]Service, len(file.GetService()))

	for i, srv := range file.GetService() {
		path := fmt.Sprintf("6,%d", i) // 6 means service.

		methods := q.generateMethods(srv, path)

		services[i].Name = srv.GetName()
		services[i].IsDeprecated = srv.GetOptions().GetDeprecated()
		services[i].ForwardMethods,
			services[i].BackwardMethods,
			services[i].BidirectionalMethods = classificateMethods(methods)
	}

	return services
}

func (q *qrpc) generateMethods(srv *pb.ServiceDescriptorProto, srvPath string) []Method {
	methods := make([]Method, len(srv.GetMethod()))

	for i, m := range srv.GetMethod() {
		path := fmt.Sprintf("%s,2,%d", srvPath, i) // 2 means method in a service.
		comment, _ := q.generateComments(path)

		methods[i].Name = m.GetName()
		methods[i].IsDeprecated = m.GetOptions().GetDeprecated()
		methods[i].Comment = comment
		methods[i].InType = q.typeName(m.GetInputType())
		methods[i].OutType = q.typeName(m.GetOutputType())
	}

	return methods
}

func (q *qrpc) renderTemplates(b io.Writer, data templateData) error {
	if err := q.callerTmpl.Execute(b, data); err != nil {
		return fmt.Errorf("cannot execute caller template: %w", err)
	}

	if err := q.handlerTmpl.Execute(b, data); err != nil {
		return fmt.Errorf("cannot execute handler template: %w", err)
	}

	return nil
}

func (q *qrpc) generateComments(path string) (string, bool) {
	loc, ok := q.comments[path]
	if !ok {
		return "", false
	}

	w := new(bytes.Buffer)
	nl := ""

	for _, line := range strings.Split(strings.TrimSuffix(loc.GetLeadingComments(), "\n"), "\n") {
		fmt.Fprintf(w, "%s//%s", nl, line)
		nl = "\n"
	}

	return w.String(), true
}

func (q *qrpc) extractComments(file *pb.FileDescriptorProto) {
	q.comments = make(map[string]*pb.SourceCodeInfo_Location)

	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		if loc.LeadingComments == nil {
			continue
		}

		var p []string
		for _, n := range loc.Path {
			p = append(p, strconv.Itoa(int(n)))
		}

		q.comments[strings.Join(p, ",")] = loc
	}
}

func classificateMethods(mm []Method) (forwardMethods, backwardMethods, bidirectionalMethods []Method) {
	for _, m := range mm {
		switch {
		case m.IsForward():
			forwardMethods = append(forwardMethods, m)
		case m.IsBackward():
			backwardMethods = append(backwardMethods, m)
		case m.IsBidirectional():
			bidirectionalMethods = append(bidirectionalMethods, m)
		}
	}

	return
}
