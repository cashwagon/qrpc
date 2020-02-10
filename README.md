# qRPC
Asynchronous protobuf based protocol.
It behaves like [gRPC](https://grpc.io/), but the RPC calls are going through message broker queue.

This project includes qRPC generator and Go library.

## Generator

qRPC provides custom generator `protoc-gen-qrpcgo`

### Installation

To use this software, you must:
- Install the standard C++ implementation of protocol buffers from
	https://developers.google.com/protocol-buffers/
- Of course, install the Go compiler and tools from
	https://golang.org/
  See
	https://golang.org/doc/install
  for details or, if you are using gccgo, follow the instructions at
	https://golang.org/doc/install/gccgo
- Grab the code from the repository and install the `proto` package.
  The simplest way is to run:
  ```
  go get -u github.com/cashwagon/qrpc/cmd/protoc-gen-qrpcgo
  ```
  The compiler plugin, `protoc-gen-qrpcgo`, will be installed in `$GOPATH/bin`
  unless `$GOBIN` is set. It must be in your `$PATH` for the protocol
  compiler, `protoc`, to find it.

**Note**: protoc-gen-qrpcgo must be used with [protoc-gen-go](https://github.com/golang/protobuf).

### Usage

How to use protoc-gen-go https://github.com/golang/protobuf/blob/master/README.md#using-protocol-buffers-with-go

To generate code compatible with qRPC pass the `qrpcgo_out` parameter to protoc alongs with `go_out`:

```shell
protoc --go_out=paths=source_relative:. --qrpcgo_out=. *.proto
```

## Library

The qRPC library uses abstract interfaces for drivers.
So you can write your own driver to support custom message broker and pass it qRPC.

For now the only supported drivers for:
- [Apache Kafka](https://kafka.apache.org/)

### Installation

To install this package, you need to install Go and setup your Go workspace on your computer. The simplest way to install the library is to run:

```shell
go get -u github.com/cashwagon/qrpc/pkg/qrpc
```

### Usage

See [examples](examples) directory for examples.

## Development

### Requirements

- [Docker](https://docs.docker.com)
- [Docker Compose](https://docs.docker.com/compose/install)
- git
- [dip](https://github.com/bibendi/dip)

### Usage

Before runnings examples or tests you need to setup `KAFKA_BROKERS` environment variable:

```shell
export KAFKA_BROKERS="kafka:9092"
```

#### First project setup

```shell
dip provision
```

#### Check installation

```shell
dip make lint
dip make test
```

#### Build

```shell
dip make
```

#### Generate protobuf files

```shell
dip make generate
```

#### Connect to postgresql shell

```shell
dip psql
```

#### List of supported commands
```shell
dip ls
```

#### Integrate dip into your shell
```shell
eval "$(dip console)"
```

#### Close (down) project
```shell
dip down
```

### Services

#### Kafka Manager

**Host**: http://localhost:9000

Manages local Kafka cluster.
You need to setup cluster in kafka-manager for the first time:
```
Cluster -> Add Cluster

Cluster Name: local
Cluster Zookeeper Hosts: zookeeper:2181
Kafka Version: 2.2.0
Enable JMX Polling: Yes
Poll consumer information: Yes

Save
```

#### Kafdrop

**Host**: http://localhost:9001

Allows to view and read topics from Kafka.

### VSCode

Project supports [VSCode remote containers](https://code.visualstudio.com/docs/remote/containers).

To start development with VSCode run `dip provision` for the first time.
Then just run command `Remote-Containers: Open Folder in Container...` in VSCode and select the project folder.
