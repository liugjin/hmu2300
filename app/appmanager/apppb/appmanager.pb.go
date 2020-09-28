// Code generated by protoc-gen-go. DO NOT EDIT.
// source: appmanager.proto

package apppb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// The request message containing publish info.
type NotifyRequest struct {
	Topic                string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Payload              string   `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Caller               string   `protobuf:"bytes,3,opt,name=caller,proto3" json:"caller,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyRequest) Reset()         { *m = NotifyRequest{} }
func (m *NotifyRequest) String() string { return proto.CompactTextString(m) }
func (*NotifyRequest) ProtoMessage()    {}
func (*NotifyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_65773fd0cc5fc11b, []int{0}
}

func (m *NotifyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyRequest.Unmarshal(m, b)
}
func (m *NotifyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyRequest.Marshal(b, m, deterministic)
}
func (m *NotifyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyRequest.Merge(m, src)
}
func (m *NotifyRequest) XXX_Size() int {
	return xxx_messageInfo_NotifyRequest.Size(m)
}
func (m *NotifyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyRequest proto.InternalMessageInfo

func (m *NotifyRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *NotifyRequest) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *NotifyRequest) GetCaller() string {
	if m != nil {
		return m.Caller
	}
	return ""
}

// The response message containing publish response
type NotifyReply struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyReply) Reset()         { *m = NotifyReply{} }
func (m *NotifyReply) String() string { return proto.CompactTextString(m) }
func (*NotifyReply) ProtoMessage()    {}
func (*NotifyReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_65773fd0cc5fc11b, []int{1}
}

func (m *NotifyReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyReply.Unmarshal(m, b)
}
func (m *NotifyReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyReply.Marshal(b, m, deterministic)
}
func (m *NotifyReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyReply.Merge(m, src)
}
func (m *NotifyReply) XXX_Size() int {
	return xxx_messageInfo_NotifyReply.Size(m)
}
func (m *NotifyReply) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyReply.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyReply proto.InternalMessageInfo

func (m *NotifyReply) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *NotifyReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*NotifyRequest)(nil), "apppb.NotifyRequest")
	proto.RegisterType((*NotifyReply)(nil), "apppb.NotifyReply")
}

func init() { proto.RegisterFile("appmanager.proto", fileDescriptor_65773fd0cc5fc11b) }

var fileDescriptor_65773fd0cc5fc11b = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0x31, 0x0f, 0x82, 0x30,
	0x10, 0x85, 0x45, 0x02, 0xc6, 0x33, 0x26, 0xe6, 0x42, 0x4c, 0xe3, 0x64, 0x98, 0x9c, 0x18, 0xd4,
	0xcd, 0xc1, 0xf8, 0x07, 0x1c, 0x58, 0x9c, 0x0f, 0xac, 0x84, 0xa4, 0xd8, 0x93, 0x96, 0xa1, 0xff,
	0xde, 0x50, 0x60, 0xd0, 0xf1, 0x7b, 0x77, 0xf9, 0xf2, 0x1e, 0x6c, 0x88, 0xb9, 0xa1, 0x37, 0x55,
	0xb2, 0xcd, 0xb8, 0xd5, 0x56, 0x63, 0x44, 0xcc, 0x5c, 0xa4, 0x0f, 0x58, 0xdf, 0xb5, 0xad, 0x5f,
	0x2e, 0x97, 0x9f, 0x4e, 0x1a, 0x8b, 0x09, 0x44, 0x56, 0x73, 0x5d, 0x8a, 0x60, 0x1f, 0x1c, 0x96,
	0xf9, 0x00, 0x28, 0x60, 0xc1, 0xe4, 0x94, 0xa6, 0xa7, 0x98, 0xfb, 0x7c, 0x42, 0xdc, 0x42, 0x5c,
	0x92, 0x52, 0xb2, 0x15, 0xa1, 0x3f, 0x8c, 0x94, 0x5e, 0x61, 0x35, 0x89, 0x59, 0xb9, 0xfe, 0xcd,
	0x58, 0xb2, 0x9d, 0xf1, 0xde, 0x28, 0x1f, 0xa9, 0x17, 0x37, 0xd2, 0x18, 0xaa, 0xe4, 0x24, 0x1e,
	0xf1, 0x78, 0x81, 0xf0, 0xc6, 0x8c, 0x67, 0x88, 0x07, 0x0f, 0x26, 0x99, 0xaf, 0x9c, 0xfd, 0xf4,
	0xdd, 0xe1, 0x5f, 0xca, 0xca, 0xa5, 0xb3, 0x22, 0xf6, 0x23, 0x4f, 0xdf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xe5, 0x47, 0x00, 0x4a, 0xf8, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AppClient is the client API for App service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AppClient interface {
	// Notify notify
	Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error)
}

type appClient struct {
	cc *grpc.ClientConn
}

func NewAppClient(cc *grpc.ClientConn) AppClient {
	return &appClient{cc}
}

func (c *appClient) Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*NotifyReply, error) {
	out := new(NotifyReply)
	err := c.cc.Invoke(ctx, "/apppb.App/Notify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppServer is the server API for App service.
type AppServer interface {
	// Notify notify
	Notify(context.Context, *NotifyRequest) (*NotifyReply, error)
}

// UnimplementedAppServer can be embedded to have forward compatible implementations.
type UnimplementedAppServer struct {
}

func (*UnimplementedAppServer) Notify(ctx context.Context, req *NotifyRequest) (*NotifyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}

func RegisterAppServer(s *grpc.Server, srv AppServer) {
	s.RegisterService(&_App_serviceDesc, srv)
}

func _App_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apppb.App/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Notify(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _App_serviceDesc = grpc.ServiceDesc{
	ServiceName: "apppb.App",
	HandlerType: (*AppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Notify",
			Handler:    _App_Notify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "appmanager.proto",
}
