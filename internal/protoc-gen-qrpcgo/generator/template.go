package generator

import (
	"fmt"
	"strings"
	"text/template"
)

type templateData struct {
	SourceFile           string
	Type                 string
	Imports              []templateImport
	GeneratedCodeVersion uint
	Services             []templateService
}

type templateImport struct {
	Alias   string
	Package string
}

type templateService struct {
	Name    string
	Methods []templateMethod
}

type templateMethod struct {
	Name    string
	InType  string
	OutType string
}

var templateFuncs = template.FuncMap{ // nolint:gochecknoglobals // there is no race
	"sprintf":             fmt.Sprintf,
	"unexport":            unexport,
	"isUnaryMethod":       isUnaryMethod,
	"isBinaryMethod":      isBinaryMethod,
	"filterBinaryMethods": filterBinaryMethods,
	"isAnyBinaryMethods":  isAnyBinaryMethods,
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func isUnaryMethod(m templateMethod) bool {
	return m.OutType == "empty.Empty"
}

func isBinaryMethod(m templateMethod) bool {
	return m.OutType != "empty.Empty"
}

func filterBinaryMethods(mm []templateMethod) []templateMethod {
	res := []templateMethod{}

	for _, m := range mm {
		if isBinaryMethod(m) {
			res = append(res, m)
		}
	}

	return res
}

func isAnyBinaryMethods(mm []templateMethod) bool {
	for _, m := range mm {
		if isBinaryMethod(m) {
			return true
		}
	}

	return false
}
