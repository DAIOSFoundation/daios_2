// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto3/proto3.proto

package proto3

import (
	fmt "fmt"
	proto "github.com/dai/go-ipfs/gxlibs/github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request_Flavour int32

const (
	Request_SWEET         Request_Flavour = 0
	Request_SOUR          Request_Flavour = 1
	Request_UMAMI         Request_Flavour = 2
	Request_GOPHERLICIOUS Request_Flavour = 3
)

var Request_Flavour_name = map[int32]string{
	0: "SWEET",
	1: "SOUR",
	2: "UMAMI",
	3: "GOPHERLICIOUS",
}

var Request_Flavour_value = map[string]int32{
	"SWEET":         0,
	"SOUR":          1,
	"UMAMI":         2,
	"GOPHERLICIOUS": 3,
}

func (x Request_Flavour) String() string {
	return proto.EnumName(Request_Flavour_name, int32(x))
}

func (Request_Flavour) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab04eb4084a521db, []int{0, 0}
}

type Request struct {
	Name                 string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Key                  []int64         `protobuf:"varint,2,rep,packed,name=key,proto3" json:"key,omitempty"`
	Taste                Request_Flavour `protobuf:"varint,3,opt,name=taste,proto3,enum=proto3.Request_Flavour" json:"taste,omitempty"`
	Book                 *Book           `protobuf:"bytes,4,opt,name=book,proto3" json:"book,omitempty"`
	Unpacked             []int64         `protobuf:"varint,5,rep,name=unpacked,proto3" json:"unpacked,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab04eb4084a521db, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Request) GetKey() []int64 {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Request) GetTaste() Request_Flavour {
	if m != nil {
		return m.Taste
	}
	return Request_SWEET
}

func (m *Request) GetBook() *Book {
	if m != nil {
		return m.Book
	}
	return nil
}

func (m *Request) GetUnpacked() []int64 {
	if m != nil {
		return m.Unpacked
	}
	return nil
}

type Book struct {
	Title                string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	RawData              []byte   `protobuf:"bytes,2,opt,name=raw_data,json=rawData,proto3" json:"raw_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Book) Reset()         { *m = Book{} }
func (m *Book) String() string { return proto.CompactTextString(m) }
func (*Book) ProtoMessage()    {}
func (*Book) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab04eb4084a521db, []int{1}
}

func (m *Book) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Book.Unmarshal(m, b)
}
func (m *Book) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Book.Marshal(b, m, deterministic)
}
func (m *Book) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Book.Merge(m, src)
}
func (m *Book) XXX_Size() int {
	return xxx_messageInfo_Book.Size(m)
}
func (m *Book) XXX_DiscardUnknown() {
	xxx_messageInfo_Book.DiscardUnknown(m)
}

var xxx_messageInfo_Book proto.InternalMessageInfo

func (m *Book) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Book) GetRawData() []byte {
	if m != nil {
		return m.RawData
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto3.Request_Flavour", Request_Flavour_name, Request_Flavour_value)
	proto.RegisterType((*Request)(nil), "proto3.Request")
	proto.RegisterType((*Book)(nil), "proto3.Book")
}

func init() { proto.RegisterFile("proto3/proto3.proto", fileDescriptor_ab04eb4084a521db) }

var fileDescriptor_ab04eb4084a521db = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x90, 0xcf, 0x4e, 0xf2, 0x40,
	0x14, 0xc5, 0x99, 0xfe, 0xf9, 0x80, 0xfb, 0xa1, 0x19, 0xaf, 0x26, 0x8e, 0x1b, 0x33, 0x61, 0xd5,
	0x0d, 0x25, 0xc1, 0x85, 0xc6, 0xb8, 0x11, 0x45, 0x25, 0x91, 0x60, 0x06, 0x89, 0x89, 0x1b, 0x33,
	0x85, 0xb1, 0x92, 0x42, 0x07, 0xcb, 0x54, 0xe2, 0xcb, 0xfa, 0x2c, 0xa6, 0x9d, 0xe2, 0xea, 0x9e,
	0x7b, 0xe7, 0xe4, 0x77, 0x32, 0x07, 0x0e, 0xd7, 0x99, 0x36, 0xfa, 0xac, 0x6b, 0x47, 0x58, 0x0e,
	0xfc, 0x67, 0xb7, 0xf6, 0x0f, 0x81, 0xba, 0x50, 0x9f, 0xb9, 0xda, 0x18, 0x44, 0xf0, 0x52, 0xb9,
	0x52, 0x8c, 0x70, 0x12, 0x34, 0x45, 0xa9, 0x91, 0x82, 0x9b, 0xa8, 0x6f, 0xe6, 0x70, 0x37, 0x70,
	0x45, 0x21, 0xb1, 0x03, 0xbe, 0x91, 0x1b, 0xa3, 0x98, 0xcb, 0x49, 0xb0, 0xdf, 0x3b, 0x0e, 0x2b,
	0x6e, 0x45, 0x09, 0xef, 0x96, 0xf2, 0x4b, 0xe7, 0x99, 0xb0, 0x2e, 0xe4, 0xe0, 0x45, 0x5a, 0x27,
	0xcc, 0xe3, 0x24, 0xf8, 0xdf, 0x6b, 0xed, 0xdc, 0x7d, 0xad, 0x13, 0x51, 0xbe, 0xe0, 0x29, 0x34,
	0xf2, 0x74, 0x2d, 0x67, 0x89, 0x9a, 0x33, 0xbf, 0xc8, 0xe9, 0x3b, 0xb4, 0x26, 0xfe, 0x6e, 0xed,
	0x2b, 0xa8, 0x57, 0x4c, 0x6c, 0x82, 0x3f, 0x79, 0x19, 0x0c, 0x9e, 0x69, 0x0d, 0x1b, 0xe0, 0x4d,
	0xc6, 0x53, 0x41, 0x49, 0x71, 0x9c, 0x8e, 0xae, 0x47, 0x43, 0xea, 0xe0, 0x01, 0xec, 0xdd, 0x8f,
	0x9f, 0x1e, 0x06, 0xe2, 0x71, 0x78, 0x33, 0x1c, 0x4f, 0x27, 0xd4, 0x6d, 0x9f, 0x83, 0x57, 0x64,
	0xe1, 0x11, 0xf8, 0x66, 0x61, 0x96, 0xbb, 0xdf, 0xd9, 0x05, 0x4f, 0xa0, 0x91, 0xc9, 0xed, 0xdb,
	0x5c, 0x1a, 0xc9, 0x1c, 0x4e, 0x82, 0x96, 0xa8, 0x67, 0x72, 0x7b, 0x2b, 0x8d, 0xec, 0x5f, 0xbe,
	0x5e, 0xc4, 0x0b, 0xf3, 0x91, 0x47, 0xe1, 0x4c, 0xaf, 0xba, 0xb1, 0x5e, 0xca, 0x34, 0xb6, 0x1d,
	0x46, 0xf9, 0xbb, 0x15, 0xb3, 0x4e, 0xac, 0xd2, 0x4e, 0xac, 0xbb, 0x46, 0x6d, 0x4c, 0xc1, 0xa8,
	0x3a, 0x8e, 0xaa, 0x76, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xec, 0x71, 0xee, 0xdb, 0x7b, 0x01,
	0x00, 0x00,
}
