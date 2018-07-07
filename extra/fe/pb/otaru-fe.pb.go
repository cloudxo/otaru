// Code generated by protoc-gen-go. DO NOT EDIT.
// source: otaru-fe.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ListHostsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListHostsRequest) Reset()         { *m = ListHostsRequest{} }
func (m *ListHostsRequest) String() string { return proto.CompactTextString(m) }
func (*ListHostsRequest) ProtoMessage()    {}
func (*ListHostsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_otaru_fe_6c9c50afbd972c7c, []int{0}
}
func (m *ListHostsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListHostsRequest.Unmarshal(m, b)
}
func (m *ListHostsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListHostsRequest.Marshal(b, m, deterministic)
}
func (dst *ListHostsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListHostsRequest.Merge(dst, src)
}
func (m *ListHostsRequest) XXX_Size() int {
	return xxx_messageInfo_ListHostsRequest.Size(m)
}
func (m *ListHostsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListHostsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListHostsRequest proto.InternalMessageInfo

type ListHostsResponse struct {
	Host                 []string `protobuf:"bytes,1,rep,name=host,proto3" json:"host,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListHostsResponse) Reset()         { *m = ListHostsResponse{} }
func (m *ListHostsResponse) String() string { return proto.CompactTextString(m) }
func (*ListHostsResponse) ProtoMessage()    {}
func (*ListHostsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_otaru_fe_6c9c50afbd972c7c, []int{1}
}
func (m *ListHostsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListHostsResponse.Unmarshal(m, b)
}
func (m *ListHostsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListHostsResponse.Marshal(b, m, deterministic)
}
func (dst *ListHostsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListHostsResponse.Merge(dst, src)
}
func (m *ListHostsResponse) XXX_Size() int {
	return xxx_messageInfo_ListHostsResponse.Size(m)
}
func (m *ListHostsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListHostsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListHostsResponse proto.InternalMessageInfo

func (m *ListHostsResponse) GetHost() []string {
	if m != nil {
		return m.Host
	}
	return nil
}

func init() {
	proto.RegisterType((*ListHostsRequest)(nil), "pb.ListHostsRequest")
	proto.RegisterType((*ListHostsResponse)(nil), "pb.ListHostsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FeServiceClient is the client API for FeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FeServiceClient interface {
	ListHosts(ctx context.Context, in *ListHostsRequest, opts ...grpc.CallOption) (*ListHostsResponse, error)
}

type feServiceClient struct {
	cc *grpc.ClientConn
}

func NewFeServiceClient(cc *grpc.ClientConn) FeServiceClient {
	return &feServiceClient{cc}
}

func (c *feServiceClient) ListHosts(ctx context.Context, in *ListHostsRequest, opts ...grpc.CallOption) (*ListHostsResponse, error) {
	out := new(ListHostsResponse)
	err := c.cc.Invoke(ctx, "/pb.FeService/ListHosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeServiceServer is the server API for FeService service.
type FeServiceServer interface {
	ListHosts(context.Context, *ListHostsRequest) (*ListHostsResponse, error)
}

func RegisterFeServiceServer(s *grpc.Server, srv FeServiceServer) {
	s.RegisterService(&_FeService_serviceDesc, srv)
}

func _FeService_ListHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListHostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeServiceServer).ListHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.FeService/ListHosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeServiceServer).ListHosts(ctx, req.(*ListHostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _FeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.FeService",
	HandlerType: (*FeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListHosts",
			Handler:    _FeService_ListHosts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "otaru-fe.proto",
}

func init() { proto.RegisterFile("otaru-fe.proto", fileDescriptor_otaru_fe_6c9c50afbd972c7c) }

var fileDescriptor_otaru_fe_6c9c50afbd972c7c = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xd9, 0x56, 0x84, 0x06, 0x95, 0x1a, 0x14, 0xca, 0x22, 0x52, 0xf6, 0xa2, 0x88, 0xdd,
	0xd8, 0x7a, 0xeb, 0xc9, 0x2a, 0x88, 0x82, 0xa0, 0xd4, 0x8b, 0x37, 0xc9, 0x2e, 0xd3, 0x6c, 0x4a,
	0xcd, 0xc4, 0x9d, 0xd9, 0xaa, 0x47, 0x05, 0x5f, 0x40, 0x1f, 0xcd, 0x57, 0xf0, 0x41, 0xa4, 0xa9,
	0x88, 0xb4, 0xa7, 0x24, 0x5f, 0xfe, 0xf9, 0x61, 0x3e, 0xb1, 0x81, 0xac, 0xcb, 0xaa, 0x33, 0x82,
	0xd4, 0x97, 0xc8, 0x28, 0x6b, 0x3e, 0x8b, 0x77, 0x0c, 0xa2, 0x99, 0x80, 0xd2, 0xde, 0x2a, 0xed,
	0x1c, 0xb2, 0x66, 0x8b, 0x8e, 0xe6, 0x89, 0xf8, 0x30, 0x1c, 0x79, 0xc7, 0x80, 0xeb, 0xd0, 0x93,
	0x36, 0x06, 0x4a, 0x85, 0x3e, 0x24, 0x96, 0xd3, 0x89, 0x14, 0xcd, 0x2b, 0x4b, 0x7c, 0x81, 0xc4,
	0x34, 0x84, 0xc7, 0x0a, 0x88, 0x93, 0x3d, 0xb1, 0xf9, 0x8f, 0x91, 0x47, 0x47, 0x20, 0xa5, 0x58,
	0x29, 0x90, 0xb8, 0x15, 0xb5, 0xeb, 0xfb, 0x8d, 0x61, 0xb8, 0xf7, 0xee, 0x45, 0xe3, 0x1c, 0x6e,
	0xa1, 0x9c, 0xda, 0x1c, 0xe4, 0x50, 0x34, 0xfe, 0xa6, 0xe4, 0x56, 0xea, 0xb3, 0x74, 0xb1, 0x38,
	0xde, 0x5e, 0xa0, 0xf3, 0xea, 0xa4, 0xf5, 0xf6, 0xf5, 0xfd, 0x59, 0x93, 0xb2, 0x19, 0x36, 0x9a,
	0x76, 0xd5, 0x08, 0xd4, 0xac, 0x9f, 0x4e, 0xdf, 0xa3, 0x8f, 0xc1, 0x6b, 0x24, 0xef, 0xc4, 0xda,
	0xf5, 0xaf, 0x86, 0xf6, 0xe0, 0xe6, 0x32, 0x39, 0x13, 0xeb, 0xe1, 0xdd, 0xf6, 0x25, 0x8e, 0x21,
	0x67, 0xb9, 0x5b, 0x30, 0x7b, 0xea, 0x2b, 0x65, 0x2c, 0x17, 0x55, 0x96, 0xe6, 0xf8, 0xa0, 0xdc,
	0x8b, 0x7e, 0x66, 0x15, 0xf4, 0xc5, 0xb2, 0x02, 0x87, 0x27, 0x81, 0x10, 0x83, 0x9f, 0xfd, 0xf7,
	0xea, 0xdd, 0xf4, 0xe8, 0x20, 0xaa, 0xf5, 0x9a, 0xda, 0xfb, 0x89, 0xcd, 0x83, 0x15, 0x35, 0x26,
	0x74, 0xfd, 0x25, 0x92, 0xad, 0x06, 0x59, 0xc7, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x61, 0x51,
	0xce, 0x32, 0x8e, 0x01, 0x00, 0x00,
}