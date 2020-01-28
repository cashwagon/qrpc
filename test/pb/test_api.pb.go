// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test_api.proto

package pb

import (
	context "context"
	fmt "fmt"
	qrpc "github.com/NightWolf007/qrpc/pkg/qrpc"
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

// FirstMethodRequest - request message for FirstMethod rpc.
type FirstMethodRequest struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FirstMethodRequest) Reset()         { *m = FirstMethodRequest{} }
func (m *FirstMethodRequest) String() string { return proto.CompactTextString(m) }
func (*FirstMethodRequest) ProtoMessage()    {}
func (*FirstMethodRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77683351be7bc655, []int{0}
}

func (m *FirstMethodRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirstMethodRequest.Unmarshal(m, b)
}
func (m *FirstMethodRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirstMethodRequest.Marshal(b, m, deterministic)
}
func (m *FirstMethodRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirstMethodRequest.Merge(m, src)
}
func (m *FirstMethodRequest) XXX_Size() int {
	return xxx_messageInfo_FirstMethodRequest.Size(m)
}
func (m *FirstMethodRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FirstMethodRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FirstMethodRequest proto.InternalMessageInfo

func (m *FirstMethodRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

// SecondMethodRequest - request message for SecondMethod rpc.
type SecondMethodRequest struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SecondMethodRequest) Reset()         { *m = SecondMethodRequest{} }
func (m *SecondMethodRequest) String() string { return proto.CompactTextString(m) }
func (*SecondMethodRequest) ProtoMessage()    {}
func (*SecondMethodRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77683351be7bc655, []int{1}
}

func (m *SecondMethodRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SecondMethodRequest.Unmarshal(m, b)
}
func (m *SecondMethodRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SecondMethodRequest.Marshal(b, m, deterministic)
}
func (m *SecondMethodRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SecondMethodRequest.Merge(m, src)
}
func (m *SecondMethodRequest) XXX_Size() int {
	return xxx_messageInfo_SecondMethodRequest.Size(m)
}
func (m *SecondMethodRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SecondMethodRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SecondMethodRequest proto.InternalMessageInfo

func (m *SecondMethodRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func init() {
	proto.RegisterType((*FirstMethodRequest)(nil), "qrpc.test.api.FirstMethodRequest")
	proto.RegisterType((*SecondMethodRequest)(nil), "qrpc.test.api.SecondMethodRequest")
}

func init() { proto.RegisterFile("test_api.proto", fileDescriptor_77683351be7bc655) }

var fileDescriptor_77683351be7bc655 = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x49, 0x2d, 0x2e,
	0x89, 0x4f, 0x2c, 0xc8, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2d, 0x2c, 0x2a, 0x48,
	0xd6, 0x03, 0x09, 0xea, 0x25, 0x16, 0x64, 0x4a, 0x49, 0xa7, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea,
	0x83, 0x25, 0x93, 0x4a, 0xd3, 0xf4, 0x53, 0x73, 0x0b, 0x4a, 0x2a, 0x21, 0x6a, 0x95, 0xd4, 0xb8,
	0x84, 0xdc, 0x32, 0x8b, 0x8a, 0x4b, 0x7c, 0x53, 0x4b, 0x32, 0xf2, 0x53, 0x82, 0x52, 0x0b, 0x4b,
	0x53, 0x8b, 0x4b, 0x84, 0x04, 0xb8, 0x98, 0x4b, 0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38,
	0x83, 0x40, 0x4c, 0x25, 0x75, 0x2e, 0xe1, 0xe0, 0xd4, 0xe4, 0xfc, 0xbc, 0x14, 0x02, 0x0a, 0x8d,
	0xe6, 0x33, 0x72, 0xb1, 0x87, 0xa4, 0x16, 0x97, 0x38, 0x06, 0x78, 0x0a, 0x79, 0x70, 0x71, 0x23,
	0x19, 0x2e, 0xa4, 0xa8, 0x87, 0xe2, 0x30, 0x3d, 0x4c, 0x8b, 0xa5, 0xc4, 0xf4, 0x20, 0x8e, 0xd5,
	0x83, 0x39, 0x56, 0xcf, 0x15, 0xe4, 0x58, 0x21, 0x2f, 0x2e, 0x1e, 0x64, 0xeb, 0x85, 0x94, 0xd0,
	0x8c, 0xc2, 0xe2, 0x36, 0x5c, 0x66, 0x39, 0xf1, 0x46, 0x71, 0x83, 0xf4, 0xe9, 0x17, 0xa5, 0x16,
	0x97, 0xe6, 0x94, 0x24, 0xb1, 0x81, 0xa5, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x4b,
	0x9a, 0x3a, 0x46, 0x01, 0x00, 0x00,
}

// This is a compile-time assertion to ensure that this generated file
// is compatible with the qrpc package it is being compiled against.
const _ = qrpc.SupportPackageIsVersion1

// TestAPIClient is the client API for TestAPI service.
type TestAPIClient interface {
	// FirstMethod tests first method.
	FirstMethod(ctx context.Context, in *FirstMethodRequest) error
	// SecondMethod tests second method.
	SecondMethod(ctx context.Context, in *SecondMethodRequest) error
}

type testAPIClient struct {
	cc *qrpc.ClientConn
}

func NewTestAPIClient(cc *qrpc.ClientConn) TestAPIClient {
	cc.SetService("qrpc.test.api.TestAPI")
	return &testAPIClient{cc}
}

func (c *testAPIClient) FirstMethod(ctx context.Context, in *FirstMethodRequest) error {
	data, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	return c.cc.Invoke(ctx, qrpc.Message{
		Method: "FirstMethod",
		Data:   data,
	})
}

func (c *testAPIClient) SecondMethod(ctx context.Context, in *SecondMethodRequest) error {
	data, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	return c.cc.Invoke(ctx, qrpc.Message{
		Method: "SecondMethod",
		Data:   data,
	})
}

// TestAPIServer is the server API for TestAPI service.
type TestAPIServer interface {
	// FirstMethod tests first method.
	FirstMethod(context.Context, *FirstMethodRequest) error
	// SecondMethod tests second method.
	SecondMethod(context.Context, *SecondMethodRequest) error
}

func RegisterTestAPIServer(s *qrpc.Server, srv TestAPIServer) {
	s.RegisterService(&_TestAPI_serviceDesc, srv)
}

func _TestAPI_FirstMethod_Handler(srv interface{}, ctx context.Context, msg []byte) error {
	in := new(FirstMethodRequest)
	if err := proto.Unmarshal(msg, in); err != nil {
		return err
	}
	return srv.(TestAPIServer).FirstMethod(ctx, in)
}

func _TestAPI_SecondMethod_Handler(srv interface{}, ctx context.Context, msg []byte) error {
	in := new(SecondMethodRequest)
	if err := proto.Unmarshal(msg, in); err != nil {
		return err
	}
	return srv.(TestAPIServer).SecondMethod(ctx, in)
}

var _TestAPI_serviceDesc = qrpc.ServiceDesc{
	ServiceName: "qrpc.test.api.TestAPI",
	HandlerType: (*TestAPIServer)(nil),
	Methods: []qrpc.MethodDesc{
		{
			MethodName: "FirstMethod",
			Handler:    _TestAPI_FirstMethod_Handler,
		},
		{
			MethodName: "SecondMethod",
			Handler:    _TestAPI_SecondMethod_Handler,
		},
	},
}
