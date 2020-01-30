// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

package pb

import (
	context "context"
	fmt "fmt"
	qrpc "github.com/cashwagon/qrpc/pkg/qrpc"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EchoRequest struct {
	Greeting             string   `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{0}
}

func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoRequest.Unmarshal(m, b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
}
func (m *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(m, src)
}
func (m *EchoRequest) XXX_Size() int {
	return xxx_messageInfo_EchoRequest.Size(m)
}
func (m *EchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EchoRequest proto.InternalMessageInfo

func (m *EchoRequest) GetGreeting() string {
	if m != nil {
		return m.Greeting
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "qrpc.example.api.EchoRequest")
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor_08134aea513e0001) }

var fileDescriptor_08134aea513e0001 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4d, 0xce, 0xc8,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x28, 0x2c, 0x2a, 0x48, 0xd6, 0x4b, 0xad, 0x48,
	0xcc, 0x2d, 0xc8, 0x49, 0xd5, 0x4b, 0x2c, 0xc8, 0x94, 0x92, 0x4e, 0xcf, 0xcf, 0x4f, 0xcf, 0x49,
	0xd5, 0x07, 0xcb, 0x27, 0x95, 0xa6, 0xe9, 0xa7, 0xe6, 0x16, 0x94, 0x54, 0x42, 0x94, 0x2b, 0x69,
	0x72, 0x71, 0xbb, 0x26, 0x67, 0xe4, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x71,
	0x71, 0xa4, 0x17, 0xa5, 0xa6, 0x96, 0x64, 0xe6, 0xa5, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06,
	0xc1, 0xf9, 0x46, 0x1e, 0x5c, 0xec, 0x20, 0xa5, 0x8e, 0x01, 0x9e, 0x42, 0xb6, 0x5c, 0x2c, 0x20,
	0xa6, 0x90, 0xac, 0x1e, 0xba, 0x6d, 0x7a, 0x48, 0xa6, 0x49, 0x89, 0xe9, 0x41, 0xac, 0xd6, 0x83,
	0x59, 0xad, 0xe7, 0x0a, 0xb2, 0xda, 0x49, 0x2b, 0x4a, 0x23, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49,
	0x2f, 0x39, 0x3f, 0x57, 0xdf, 0x2f, 0x33, 0x3d, 0xa3, 0x24, 0x3c, 0x3f, 0x27, 0xcd, 0xc0, 0xc0,
	0x5c, 0x1f, 0x64, 0x9e, 0x3e, 0xd4, 0xbc, 0x62, 0xfd, 0x82, 0xa4, 0x24, 0x36, 0xb0, 0x5e, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x68, 0x6a, 0xf5, 0x78, 0xe4, 0x00, 0x00, 0x00,
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the qrpc package it is being compiled against.
const _ = qrpc.SupportPackageIsVersion1

// EchoAPIClient is the client API for EchoAPI service.
type EchoAPIClient interface {
	Echo(ctx context.Context, in *EchoRequest) error
}

type echoAPIClient struct {
	cc *qrpc.ClientConn
}

func NewEchoAPIClient(cc *qrpc.ClientConn) EchoAPIClient {
	cc.SetService("qrpc.example.api.EchoAPI")
	return &echoAPIClient{cc}
}

func (c *echoAPIClient) Echo(ctx context.Context, in *EchoRequest) error {
	data, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	return c.cc.Invoke(ctx, qrpc.Message{
		Method: "Echo",
		Data:   data,
	})
}

// EchoAPIServer is the server API for EchoAPI service.
type EchoAPIServer interface {
	Echo(context.Context, *EchoRequest) error
}

func RegisterEchoAPIServer(s *qrpc.Server, srv EchoAPIServer) {
	s.RegisterService(&_EchoAPI_serviceDesc, srv)
}

func _EchoAPI_Echo_Handler(srv interface{}, ctx context.Context, msg []byte) error {
	in := new(EchoRequest)
	if err := proto.Unmarshal(msg, in); err != nil {
		return err
	}
	return srv.(EchoAPIServer).Echo(ctx, in)
}

var _EchoAPI_serviceDesc = qrpc.ServiceDesc{
	ServiceName: "qrpc.example.api.EchoAPI",
	HandlerType: (*EchoAPIServer)(nil),
	Methods: []qrpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoAPI_Echo_Handler,
		},
	},
}
