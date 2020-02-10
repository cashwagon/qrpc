// Package generator provides generator to generate qRPC files from proto files
package generator

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

const generatedCodeVersion = 1

const selfAlias = "_pb"

var isGoKeyword = map[string]bool{ // nolint:gochecknoglobals // there is no racing
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"else":        true,
	"defer":       true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// Import represents simple go import structure
type Import struct {
	Alias   string
	Package string
}

// Generator represents the main generator struct
// It is used to generate files for qrpc from proto
type Generator struct {
	Request     *plugin.CodeGeneratorRequest
	Response    *plugin.CodeGeneratorResponse
	file        *descriptor.FileDescriptorProto
	filesByName map[string]*descriptor.FileDescriptorProto
	imports     map[string]Import
}

// New creates new generator
func New() *Generator {
	return &Generator{
		Request:     new(plugin.CodeGeneratorRequest),
		Response:    new(plugin.CodeGeneratorResponse),
		filesByName: make(map[string]*descriptor.FileDescriptorProto),
		imports:     make(map[string]Import),
	}
}

// Generate generates all qrpc files from all proto files
func (g *Generator) Generate() error {
	for _, file := range g.Request.GetProtoFile() {
		g.filesByName[file.GetName()] = file

		pkg := file.GetOptions().GetGoPackage()
		g.imports[file.GetPackage()] = Import{
			Alias:   pkgAlias(pkg),
			Package: pkg,
		}
	}

	for _, file := range g.Request.GetProtoFile() {
		g.file = file
		if err := g.generate(); err != nil {
			return err
		}
	}

	return nil
}

// generate generates qRPC files from one proto file
func (g *Generator) generate() error {
	if len(g.file.GetService()) == 0 {
		return nil
	}

	data := templateData{
		SourceFile:           g.file.GetName(),
		Imports:              g.generateImports(),
		GeneratedCodeVersion: generatedCodeVersion,
		Services:             g.generateServices(),
	}

	err := g.generateFile(&data, callerTemplate, "caller")
	if err != nil {
		return err
	}

	err = g.generateFile(&data, handlerTemplate, "handler")
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateFile(data *templateData, tmpl, tp string) error {
	t, err := template.New(g.file.GetName()).Funcs(templateFuncs).Parse(tmpl)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}

	if err := t.Execute(&b, *data); err != nil {
		return err
	}

	g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(goFileName(tp, g.file.GetName())),
		Content: proto.String(b.String()),
	})

	return nil
}

func (g *Generator) generateImports() []templateImport {
	imports := []templateImport{}

	for _, d := range g.file.GetDependency() {
		f := g.filesByName[d]
		imp := g.imports[f.GetPackage()]

		// Skip own package
		if imp.Package == g.file.GetOptions().GetGoPackage() {
			continue
		}

		imports = append(imports, templateImport(imp))
	}

	return append(imports, templateImport{
		Alias:   selfAlias,
		Package: g.file.GetOptions().GetGoPackage(),
	})
}

func (g *Generator) generateServices() []templateService {
	services := make([]templateService, len(g.file.GetService()))

	for i, s := range g.file.GetService() {
		services[i].Name = s.GetName()
		services[i].Methods = g.generateMethods(s)
	}

	return services
}

func (g *Generator) generateMethods(s *descriptor.ServiceDescriptorProto) []templateMethod {
	methods := make([]templateMethod, len(s.GetMethod()))
	for i, m := range s.GetMethod() {
		methods[i].Name = m.GetName()
		methods[i].InType = g.getType(m.GetInputType())
		methods[i].OutType = g.getType(m.GetOutputType())
	}

	return methods
}

func (g *Generator) getType(tp string) string {
	var pkgName, tpName string
	if i := strings.LastIndex(tp, "."); i >= 0 {
		pkgName = tp[0:i]
		if pkgName[0] == '.' {
			pkgName = pkgName[1:]
		}

		tpName = tp[i+1:]
	}

	imp, ok := g.imports[pkgName]
	if ok {
		if imp.Package == g.file.GetOptions().GetGoPackage() {
			return fmt.Sprintf("%s.%s", selfAlias, tpName)
		}

		return fmt.Sprintf("%s.%s", imp.Alias, tpName)
	}

	return tpName
}

func goFileName(tp, name string) string {
	if ext := path.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}

	return fmt.Sprintf("%s/%s.pb.go", tp, name)
}

// pkgName returns alias for package path
func pkgAlias(pkg string) string {
	alias := pkgBaseName(pkg)

	// Replace all unknown symbols to underscore
	alias = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, alias)

	// Identifier must not be keyword or predeclared identifier: insert _.
	if isGoKeyword[alias] {
		alias = "_" + alias
	}

	// Identifier must not begin with digit: insert _.
	if r, _ := utf8.DecodeRuneInString(alias); unicode.IsDigit(r) {
		alias = "_" + alias
	}

	return alias
}

// pkgBaseName returns the last path element of the name, with the last dotted suffix removed.
func pkgBaseName(name string) string {
	// First, find the last element
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	// Now drop the suffix
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[0:i]
	}

	return name
}
