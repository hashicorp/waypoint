// Code generated by protoc-gen-go. DO NOT EDIT.
// source: waypoint/builtin/files/plugin.proto

package files

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Files struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Files) Reset()         { *m = Files{} }
func (m *Files) String() string { return proto.CompactTextString(m) }
func (*Files) ProtoMessage()    {}
func (*Files) Descriptor() ([]byte, []int) {
	return fileDescriptor_880353e3a77dc7cc, []int{0}
}

func (m *Files) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Files.Unmarshal(m, b)
}
func (m *Files) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Files.Marshal(b, m, deterministic)
}
func (m *Files) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Files.Merge(m, src)
}
func (m *Files) XXX_Size() int {
	return xxx_messageInfo_Files.Size(m)
}
func (m *Files) XXX_DiscardUnknown() {
	xxx_messageInfo_Files.DiscardUnknown(m)
}

var xxx_messageInfo_Files proto.InternalMessageInfo

func (m *Files) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func init() {
	proto.RegisterType((*Files)(nil), "files.Files")
}

func init() {
	proto.RegisterFile("waypoint/builtin/files/plugin.proto", fileDescriptor_880353e3a77dc7cc)
}

var fileDescriptor_880353e3a77dc7cc = []byte{
	// 100 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2e, 0x4f, 0xac, 0x2c,
	0xc8, 0xcf, 0xcc, 0x2b, 0xd1, 0x4f, 0x2a, 0xcd, 0xcc, 0x29, 0xc9, 0xcc, 0xd3, 0x4f, 0xcb, 0xcc,
	0x49, 0x2d, 0xd6, 0x2f, 0xc8, 0x29, 0x4d, 0xcf, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0x05, 0x8b, 0x29, 0x49, 0x73, 0xb1, 0xba, 0x81, 0x18, 0x42, 0x42, 0x5c, 0x2c, 0x05, 0x89,
	0x25, 0x19, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x93, 0x44, 0x94, 0x18, 0x76,
	0xa3, 0x92, 0xd8, 0xc0, 0x86, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x46, 0x1a, 0xa3, 0x2a,
	0x6b, 0x00, 0x00, 0x00,
}
