// protoc-gen-qrpcgo is a plugin for the Google protocol buffer compiler to generate
// Go code.  Run it by building this program and putting it in your path with
// the name
// 	protoc-gen-qrpcgo
// That word 'go' at the end becomes part of the option string set for the
// protocol compiler, so once the protocol compiler (protoc) is installed
// you can run
// 	protoc --qrpcgo_out=output_directory input_directory/file.proto
// to generate Go bindings for the protocol defined by file.proto.
// With that input, the output will be written to
// 	output_directory/file.pb.go
//
// The generated code is documented in the package comment for
// the library.
//
// See the README and documentation for protocol buffers to learn more:
// 	https://developers.google.com/protocol-buffers/

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/cashwagon/qrpc/internal/protoc-gen-qrpcgo/generator"
	"github.com/golang/protobuf/proto"
)

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logErr(err, "failed to read input")
	}

	if err = proto.Unmarshal(data, g.Request); err != nil {
		logErr(err, "failed to parse input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		logFail("no files to generate")
	}

	if err = g.Generate(); err != nil {
		logErr(err, "failed to generate files")
	}

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		logErr(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		logErr(err, "failed to write output proto")
	}
}

func logErr(err error, msg string) {
	log.Fatalf("protoc-gen-qrpcgo: error: %s: %v", msg, err)
}

func logFail(msg string) {
	log.Fatalf("protoc-gen-qrpcgo: error: %s", msg)
}
