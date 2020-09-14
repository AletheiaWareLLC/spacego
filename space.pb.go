// Code generated by protoc-gen-go. DO NOT EDIT.
// source: space.proto

package spacego

import (
	fmt "fmt"
	bcgo "github.com/AletheiaWareLLC/bcgo"
	financego "github.com/AletheiaWareLLC/financego"
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

type Meta struct {
	// Name of file
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Size of file in bytes
	Size uint64 `protobuf:"fixed64,2,opt,name=size,proto3" json:"size,omitempty"`
	// MIME type of file
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{0}
}

func (m *Meta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Meta.Unmarshal(m, b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return xxx_messageInfo_Meta.Size(m)
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func (m *Meta) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Meta) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Meta) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type Preview struct {
	// MIME type of preview
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Preview data
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Preview width
	Width uint32 `protobuf:"varint,3,opt,name=width,proto3" json:"width,omitempty"`
	// Preview height
	Height               uint32   `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Preview) Reset()         { *m = Preview{} }
func (m *Preview) String() string { return proto.CompactTextString(m) }
func (*Preview) ProtoMessage()    {}
func (*Preview) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{1}
}

func (m *Preview) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Preview.Unmarshal(m, b)
}
func (m *Preview) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Preview.Marshal(b, m, deterministic)
}
func (m *Preview) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Preview.Merge(m, src)
}
func (m *Preview) XXX_Size() int {
	return xxx_messageInfo_Preview.Size(m)
}
func (m *Preview) XXX_DiscardUnknown() {
	xxx_messageInfo_Preview.DiscardUnknown(m)
}

var xxx_messageInfo_Preview proto.InternalMessageInfo

func (m *Preview) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Preview) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Preview) GetWidth() uint32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Preview) GetHeight() uint32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type Share struct {
	// The reference to the meta being shared
	MetaReference *bcgo.Reference `protobuf:"bytes,1,opt,name=meta_reference,json=metaReference,proto3" json:"meta_reference,omitempty"`
	// The encryption key for meta
	MetaKey []byte `protobuf:"bytes,2,opt,name=meta_key,json=metaKey,proto3" json:"meta_key,omitempty"`
	// References to chunks being shared are in bc.Record.Reference of metadata
	// A list of encryption keys for chunks
	ChunkKey [][]byte `protobuf:"bytes,3,rep,name=chunk_key,json=chunkKey,proto3" json:"chunk_key,omitempty"`
	// A list of references to previews
	PreviewReference []*bcgo.Reference `protobuf:"bytes,4,rep,name=preview_reference,json=previewReference,proto3" json:"preview_reference,omitempty"`
	// A list of encryption keys for previews
	PreviewKey           [][]byte `protobuf:"bytes,5,rep,name=preview_key,json=previewKey,proto3" json:"preview_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Share) Reset()         { *m = Share{} }
func (m *Share) String() string { return proto.CompactTextString(m) }
func (*Share) ProtoMessage()    {}
func (*Share) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{2}
}

func (m *Share) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Share.Unmarshal(m, b)
}
func (m *Share) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Share.Marshal(b, m, deterministic)
}
func (m *Share) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Share.Merge(m, src)
}
func (m *Share) XXX_Size() int {
	return xxx_messageInfo_Share.Size(m)
}
func (m *Share) XXX_DiscardUnknown() {
	xxx_messageInfo_Share.DiscardUnknown(m)
}

var xxx_messageInfo_Share proto.InternalMessageInfo

func (m *Share) GetMetaReference() *bcgo.Reference {
	if m != nil {
		return m.MetaReference
	}
	return nil
}

func (m *Share) GetMetaKey() []byte {
	if m != nil {
		return m.MetaKey
	}
	return nil
}

func (m *Share) GetChunkKey() [][]byte {
	if m != nil {
		return m.ChunkKey
	}
	return nil
}

func (m *Share) GetPreviewReference() []*bcgo.Reference {
	if m != nil {
		return m.PreviewReference
	}
	return nil
}

func (m *Share) GetPreviewKey() [][]byte {
	if m != nil {
		return m.PreviewKey
	}
	return nil
}

type Tag struct {
	// The value of tag applied to meta
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	// The reason for tagging
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{3}
}

func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (m *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(m, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Tag) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type Miner struct {
	Merchant             *financego.Merchant `protobuf:"bytes,1,opt,name=merchant,proto3" json:"merchant,omitempty"`
	Service              *financego.Service  `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Miner) Reset()         { *m = Miner{} }
func (m *Miner) String() string { return proto.CompactTextString(m) }
func (*Miner) ProtoMessage()    {}
func (*Miner) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{4}
}

func (m *Miner) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Miner.Unmarshal(m, b)
}
func (m *Miner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Miner.Marshal(b, m, deterministic)
}
func (m *Miner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Miner.Merge(m, src)
}
func (m *Miner) XXX_Size() int {
	return xxx_messageInfo_Miner.Size(m)
}
func (m *Miner) XXX_DiscardUnknown() {
	xxx_messageInfo_Miner.DiscardUnknown(m)
}

var xxx_messageInfo_Miner proto.InternalMessageInfo

func (m *Miner) GetMerchant() *financego.Merchant {
	if m != nil {
		return m.Merchant
	}
	return nil
}

func (m *Miner) GetService() *financego.Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type Registrar struct {
	Merchant             *financego.Merchant `protobuf:"bytes,1,opt,name=merchant,proto3" json:"merchant,omitempty"`
	Service              *financego.Service  `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Registrar) Reset()         { *m = Registrar{} }
func (m *Registrar) String() string { return proto.CompactTextString(m) }
func (*Registrar) ProtoMessage()    {}
func (*Registrar) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8a3f24abfdc04ca, []int{5}
}

func (m *Registrar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Registrar.Unmarshal(m, b)
}
func (m *Registrar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Registrar.Marshal(b, m, deterministic)
}
func (m *Registrar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Registrar.Merge(m, src)
}
func (m *Registrar) XXX_Size() int {
	return xxx_messageInfo_Registrar.Size(m)
}
func (m *Registrar) XXX_DiscardUnknown() {
	xxx_messageInfo_Registrar.DiscardUnknown(m)
}

var xxx_messageInfo_Registrar proto.InternalMessageInfo

func (m *Registrar) GetMerchant() *financego.Merchant {
	if m != nil {
		return m.Merchant
	}
	return nil
}

func (m *Registrar) GetService() *financego.Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func init() {
	proto.RegisterType((*Meta)(nil), "space.Meta")
	proto.RegisterType((*Preview)(nil), "space.Preview")
	proto.RegisterType((*Share)(nil), "space.Share")
	proto.RegisterType((*Tag)(nil), "space.Tag")
	proto.RegisterType((*Miner)(nil), "space.Miner")
	proto.RegisterType((*Registrar)(nil), "space.Registrar")
}

func init() { proto.RegisterFile("space.proto", fileDescriptor_b8a3f24abfdc04ca) }

var fileDescriptor_b8a3f24abfdc04ca = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0x4b, 0x6b, 0xdb, 0x40,
	0x10, 0x46, 0x95, 0xe5, 0xc7, 0x38, 0x2e, 0xc9, 0x52, 0x82, 0x9b, 0x1e, 0x6a, 0x74, 0x32, 0x85,
	0x2a, 0x90, 0xf4, 0xd4, 0x5b, 0xdd, 0x4b, 0x21, 0x31, 0x84, 0x75, 0xa1, 0xd0, 0x4b, 0x18, 0x6d,
	0xc6, 0xda, 0x25, 0xf1, 0xca, 0xac, 0xd6, 0x36, 0xee, 0xcf, 0xec, 0x2f, 0x2a, 0xfb, 0x90, 0x5c,
	0xe8, 0xb9, 0x27, 0xcd, 0xf7, 0xd0, 0x7c, 0xbb, 0xb3, 0x03, 0xe3, 0x66, 0x8b, 0x82, 0x8a, 0xad,
	0xa9, 0x6d, 0xcd, 0x32, 0x0f, 0xae, 0x86, 0xa5, 0x08, 0xc4, 0xd5, 0x64, 0xad, 0x34, 0xea, 0x56,
	0xcf, 0x17, 0xd0, 0x5b, 0x92, 0x45, 0xc6, 0xa0, 0xa7, 0x71, 0x43, 0xd3, 0x64, 0x96, 0xcc, 0x47,
	0xdc, 0xd7, 0x8e, 0x6b, 0xd4, 0x2f, 0x9a, 0xbe, 0x9a, 0x25, 0xf3, 0x3e, 0xf7, 0xb5, 0xe3, 0xec,
	0x71, 0x4b, 0xd3, 0x34, 0xf8, 0x5c, 0x9d, 0x3f, 0xc2, 0xe0, 0xc1, 0xd0, 0x5e, 0xd1, 0xa1, 0x93,
	0x93, 0x93, 0xec, 0xb8, 0x27, 0xb4, 0xe8, 0xdb, 0x9c, 0x71, 0x5f, 0xb3, 0x37, 0x90, 0x1d, 0xd4,
	0x93, 0x95, 0xbe, 0xcf, 0x84, 0x07, 0xc0, 0x2e, 0xa1, 0x2f, 0x49, 0x55, 0xd2, 0x4e, 0x7b, 0x9e,
	0x8e, 0x28, 0xff, 0x9d, 0x40, 0xb6, 0x92, 0x68, 0x88, 0x7d, 0x82, 0xd7, 0x1b, 0xb2, 0xf8, 0x68,
	0x68, 0x4d, 0x86, 0xb4, 0x08, 0x49, 0xe3, 0x9b, 0x49, 0x51, 0x8a, 0x82, 0xb7, 0x24, 0x9f, 0x38,
	0x53, 0x07, 0xd9, 0x5b, 0x18, 0xfa, 0xbf, 0x9e, 0xe9, 0x18, 0x4f, 0x31, 0x70, 0xf8, 0x8e, 0x8e,
	0xec, 0x1d, 0x8c, 0x84, 0xdc, 0xe9, 0x67, 0xaf, 0xa5, 0xb3, 0x74, 0x7e, 0xc6, 0x87, 0x9e, 0x70,
	0xe2, 0x67, 0xb8, 0xd8, 0x86, 0x8b, 0xfd, 0x15, 0xd8, 0x9b, 0xa5, 0xff, 0x06, 0x9e, 0x47, 0xdf,
	0x29, 0xf3, 0x3d, 0x8c, 0xdb, 0x7f, 0x5d, 0xeb, 0xcc, 0xb7, 0x86, 0x48, 0xdd, 0xd1, 0x31, 0xbf,
	0x85, 0xf4, 0x3b, 0x56, 0x6e, 0x12, 0x7b, 0x7c, 0xd9, 0xb5, 0x23, 0x0b, 0xc0, 0x4d, 0xc2, 0x10,
	0x36, 0xb5, 0xf6, 0xe7, 0x1d, 0xf1, 0x88, 0xf2, 0x12, 0xb2, 0xa5, 0xd2, 0x64, 0xd8, 0x47, 0x77,
	0x25, 0x23, 0x24, 0x6a, 0x1b, 0x47, 0x70, 0x51, 0xb4, 0x2f, 0xbb, 0x8c, 0x02, 0xef, 0x2c, 0xec,
	0x03, 0x0c, 0x1a, 0x32, 0x7b, 0x25, 0xc2, 0x6b, 0x8e, 0x6f, 0xce, 0x3b, 0xf7, 0x2a, 0xf0, 0xbc,
	0x35, 0xe4, 0x6b, 0x18, 0x71, 0xaa, 0x54, 0x63, 0x0d, 0xfe, 0xcf, 0x9c, 0xc5, 0x37, 0xb8, 0x14,
	0xf5, 0xa6, 0xc0, 0x17, 0xb2, 0x92, 0x14, 0x1e, 0xd0, 0x50, 0xe1, 0xb7, 0x75, 0x01, 0x2b, 0xf7,
	0x79, 0x70, 0x0b, 0xfa, 0x33, 0xaf, 0x94, 0x95, 0xbb, 0xb2, 0x10, 0xf5, 0xe6, 0xfa, 0x4b, 0xb4,
	0xfe, 0x40, 0x43, 0xf7, 0xf7, 0x5f, 0xaf, 0xbd, 0xbb, 0xaa, 0xcb, 0xbe, 0xdf, 0xe5, 0xdb, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x24, 0x45, 0x1b, 0xb5, 0xfa, 0x02, 0x00, 0x00,
}
