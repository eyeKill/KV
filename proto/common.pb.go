// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package proto

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

type Operation int32

const (
	Operation_GET                Operation = 0
	Operation_PUT                Operation = 1
	Operation_DELETE             Operation = 2
	Operation_START_TRANSACTION  Operation = 3
	Operation_COMMIT_TRANSACTION Operation = 4
)

var Operation_name = map[int32]string{
	0: "GET",
	1: "PUT",
	2: "DELETE",
	3: "START_TRANSACTION",
	4: "COMMIT_TRANSACTION",
}

var Operation_value = map[string]int32{
	"GET":                0,
	"PUT":                1,
	"DELETE":             2,
	"START_TRANSACTION":  3,
	"COMMIT_TRANSACTION": 4,
}

func (x Operation) String() string {
	return proto.EnumName(Operation_name, int32(x))
}

func (Operation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

type Status int32

const (
	Status_OK        Status = 0
	Status_ENOENT    Status = 1
	Status_ENOSERVER Status = 2
	Status_EFAILED   Status = 3
)

var Status_name = map[int32]string{
	0: "OK",
	1: "ENOENT",
	2: "ENOSERVER",
	3: "EFAILED",
}

var Status_value = map[string]int32{
	"OK":        0,
	"ENOENT":    1,
	"ENOSERVER": 2,
	"EFAILED":   3,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

type Key struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Key) Reset()         { *m = Key{} }
func (m *Key) String() string { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()    {}
func (*Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *Key) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Key.Unmarshal(m, b)
}
func (m *Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Key.Marshal(b, m, deterministic)
}
func (m *Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Key.Merge(m, src)
}
func (m *Key) XXX_Size() int {
	return xxx_messageInfo_Key.Size(m)
}
func (m *Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Key proto.InternalMessageInfo

func (m *Key) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type Value struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type KVPair struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KVPair) Reset()         { *m = KVPair{} }
func (m *KVPair) String() string { return proto.CompactTextString(m) }
func (*KVPair) ProtoMessage()    {}
func (*KVPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{2}
}

func (m *KVPair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KVPair.Unmarshal(m, b)
}
func (m *KVPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KVPair.Marshal(b, m, deterministic)
}
func (m *KVPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KVPair.Merge(m, src)
}
func (m *KVPair) XXX_Size() int {
	return xxx_messageInfo_KVPair.Size(m)
}
func (m *KVPair) XXX_DiscardUnknown() {
	xxx_messageInfo_KVPair.DiscardUnknown(m)
}

var xxx_messageInfo_KVPair proto.InternalMessageInfo

func (m *KVPair) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KVPair) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type WorkerId struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorkerId) Reset()         { *m = WorkerId{} }
func (m *WorkerId) String() string { return proto.CompactTextString(m) }
func (*WorkerId) ProtoMessage()    {}
func (*WorkerId) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{3}
}

func (m *WorkerId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkerId.Unmarshal(m, b)
}
func (m *WorkerId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkerId.Marshal(b, m, deterministic)
}
func (m *WorkerId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkerId.Merge(m, src)
}
func (m *WorkerId) XXX_Size() int {
	return xxx_messageInfo_WorkerId.Size(m)
}
func (m *WorkerId) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkerId.DiscardUnknown(m)
}

var xxx_messageInfo_WorkerId proto.InternalMessageInfo

func (m *WorkerId) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type BackupEntry struct {
	Op                   Operation `protobuf:"varint,1,opt,name=op,proto3,enum=kv.proto.Operation" json:"op,omitempty"`
	Version              uint64    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Key                  string    `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value                string    `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *BackupEntry) Reset()         { *m = BackupEntry{} }
func (m *BackupEntry) String() string { return proto.CompactTextString(m) }
func (*BackupEntry) ProtoMessage()    {}
func (*BackupEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{4}
}

func (m *BackupEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackupEntry.Unmarshal(m, b)
}
func (m *BackupEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackupEntry.Marshal(b, m, deterministic)
}
func (m *BackupEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackupEntry.Merge(m, src)
}
func (m *BackupEntry) XXX_Size() int {
	return xxx_messageInfo_BackupEntry.Size(m)
}
func (m *BackupEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_BackupEntry.DiscardUnknown(m)
}

var xxx_messageInfo_BackupEntry proto.InternalMessageInfo

func (m *BackupEntry) GetOp() Operation {
	if m != nil {
		return m.Op
	}
	return Operation_GET
}

func (m *BackupEntry) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *BackupEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *BackupEntry) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterEnum("kv.proto.Operation", Operation_name, Operation_value)
	proto.RegisterEnum("kv.proto.Status", Status_name, Status_value)
	proto.RegisterType((*Key)(nil), "kv.proto.Key")
	proto.RegisterType((*Value)(nil), "kv.proto.Value")
	proto.RegisterType((*KVPair)(nil), "kv.proto.KVPair")
	proto.RegisterType((*WorkerId)(nil), "kv.proto.WorkerId")
	proto.RegisterType((*BackupEntry)(nil), "kv.proto.BackupEntry")
}

func init() {
	proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206)
}

var fileDescriptor_555bd8c177793206 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xcf, 0x4f, 0xc2, 0x30,
	0x14, 0xc7, 0xa1, 0x83, 0x0d, 0x1e, 0x42, 0xea, 0xf3, 0x17, 0x31, 0x31, 0x31, 0xf3, 0x62, 0x38,
	0x10, 0xa3, 0x17, 0xaf, 0x03, 0xaa, 0x59, 0x80, 0x8d, 0x74, 0x15, 0x8d, 0x17, 0x33, 0x61, 0x87,
	0x65, 0xb2, 0x2e, 0xa5, 0x90, 0xf0, 0xdf, 0x9b, 0x0d, 0x51, 0x8c, 0x9e, 0xfa, 0x7d, 0x7d, 0xf9,
	0x7c, 0xfb, 0x29, 0x1c, 0xcc, 0xe4, 0x62, 0x21, 0xd3, 0x6e, 0xa6, 0xa4, 0x96, 0x58, 0x4b, 0xd6,
	0xdb, 0x64, 0x9f, 0x81, 0x31, 0x8c, 0x36, 0x48, 0xc1, 0x48, 0xa2, 0x4d, 0xbb, 0x7c, 0x59, 0xbe,
	0xae, 0xf3, 0x3c, 0xda, 0x17, 0x50, 0x9d, 0x86, 0x1f, 0xab, 0x08, 0x8f, 0xa1, 0xba, 0xce, 0xc3,
	0xd7, 0x72, 0x3b, 0xd8, 0x37, 0x60, 0x0e, 0xa7, 0x93, 0x30, 0x56, 0x7f, 0xd1, 0x1f, 0x82, 0xec,
	0x13, 0xe7, 0x50, 0x7b, 0x96, 0x2a, 0x89, 0x94, 0x3b, 0xc7, 0x16, 0x90, 0x78, 0x5e, 0x20, 0x4d,
	0x4e, 0xe2, 0xb9, 0xad, 0xa1, 0xd1, 0x0b, 0x67, 0xc9, 0x2a, 0x63, 0xa9, 0x56, 0x1b, 0xbc, 0x02,
	0x22, 0xb3, 0x62, 0xdd, 0xba, 0x3d, 0xea, 0xee, 0x5c, 0xbb, 0x7e, 0x16, 0xa9, 0x50, 0xc7, 0x32,
	0xe5, 0x44, 0x66, 0xd8, 0x06, 0x6b, 0x1d, 0xa9, 0x65, 0x2c, 0xd3, 0xe2, 0x9d, 0x0a, 0xdf, 0x8d,
	0x3b, 0x23, 0xe3, 0x1f, 0xa3, 0xca, 0x9e, 0x51, 0xe7, 0x05, 0xea, 0xdf, 0x95, 0x68, 0x81, 0xf1,
	0xc8, 0x04, 0x2d, 0xe5, 0x61, 0xf2, 0x24, 0x68, 0x19, 0x01, 0xcc, 0x01, 0x1b, 0x31, 0xc1, 0x28,
	0xc1, 0x13, 0x38, 0x0c, 0x84, 0xc3, 0xc5, 0x9b, 0xe0, 0x8e, 0x17, 0x38, 0x7d, 0xe1, 0xfa, 0x1e,
	0x35, 0xf0, 0x14, 0xb0, 0xef, 0x8f, 0xc7, 0xee, 0xef, 0xfb, 0x4a, 0xe7, 0x1e, 0xcc, 0x40, 0x87,
	0x7a, 0xb5, 0x44, 0x13, 0x88, 0x3f, 0xa4, 0xa5, 0xbc, 0x8c, 0x79, 0x3e, 0xf3, 0xf2, 0xe2, 0x26,
	0xd4, 0x99, 0xe7, 0x07, 0x8c, 0x4f, 0x19, 0xa7, 0x04, 0x1b, 0x60, 0xb1, 0x07, 0xc7, 0x1d, 0xb1,
	0x01, 0x35, 0x7a, 0xd6, 0x6b, 0xb5, 0xf8, 0xec, 0xbb, 0x59, 0x1c, 0x77, 0x9f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x88, 0x82, 0x70, 0xa3, 0xb9, 0x01, 0x00, 0x00,
}
