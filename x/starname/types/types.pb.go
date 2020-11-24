// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: x/starname/types/types.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

// Resource defines a resource an account can resolve to
type Resource struct {
	// URI defines the ID of the resource
	URI string `protobuf:"bytes,1,opt,name=URI,proto3" json:"uri"`
	// Resource is the resource
	Resource             string   `protobuf:"bytes,2,opt,name=Resource,proto3" json:"resource"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resource) Reset()         { *m = Resource{} }
func (m *Resource) String() string { return proto.CompactTextString(m) }
func (*Resource) ProtoMessage()    {}
func (*Resource) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb68cf137d2ded7a, []int{0}
}
func (m *Resource) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Resource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Resource.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Resource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resource.Merge(m, src)
}
func (m *Resource) XXX_Size() int {
	return m.Size()
}
func (m *Resource) XXX_DiscardUnknown() {
	xxx_messageInfo_Resource.DiscardUnknown(m)
}

var xxx_messageInfo_Resource proto.InternalMessageInfo

func (m *Resource) GetURI() string {
	if m != nil {
		return m.URI
	}
	return ""
}

func (m *Resource) GetResource() string {
	if m != nil {
		return m.Resource
	}
	return ""
}

// Domain defines a domain
type Domain struct {
	// Name is the name of the domain
	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"name"`
	// Admin is the owner of the domain
	Admin  github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=Admin,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"admin"`
	Broker github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=Broker,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"broker"`
	// ValidUntil is a unix timestamp defines the time when the domain will become invalid in seconds
	ValidUntil int64 `protobuf:"varint,4,opt,name=ValidUntil,proto3" json:"valid_until"`
	// Type defines the type of the domain
	Type                 DomainType `protobuf:"bytes,5,opt,name=Type,proto3,casttype=DomainType" json:"type"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Domain) Reset()         { *m = Domain{} }
func (m *Domain) String() string { return proto.CompactTextString(m) }
func (*Domain) ProtoMessage()    {}
func (*Domain) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb68cf137d2ded7a, []int{1}
}
func (m *Domain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Domain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Domain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Domain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Domain.Merge(m, src)
}
func (m *Domain) XXX_Size() int {
	return m.Size()
}
func (m *Domain) XXX_DiscardUnknown() {
	xxx_messageInfo_Domain.DiscardUnknown(m)
}

var xxx_messageInfo_Domain proto.InternalMessageInfo

func (m *Domain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Domain) GetAdmin() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Admin
	}
	return nil
}

func (m *Domain) GetBroker() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Broker
	}
	return nil
}

func (m *Domain) GetValidUntil() int64 {
	if m != nil {
		return m.ValidUntil
	}
	return 0
}

func (m *Domain) GetType() DomainType {
	if m != nil {
		return m.Type
	}
	return ""
}

// Account defines an account that belongs to a domain
// NOTE: It should not be confused with cosmos-sdk auth account
// github.com/cosmos/cosmos-sdk/x/auth.Account
type Account struct {
	// Domain references the domain this account belongs to
	Domain string `protobuf:"bytes,1,opt,name=Domain,proto3" json:"domain"`
	// Name is the name of the account
	Name *string `protobuf:"bytes,2,opt,name=Name,proto3,wktptr" json:"name"`
	// Owner is the address that owns the account
	Owner github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=Owner,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"owner"`
	// Broker can be empty
	// it identifies an entity that facilitated the transaction of the account
	Broker github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,4,opt,name=Broker,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"broker"`
	// ValidUntil defines a unix timestamp of the expiration of the account in seconds
	ValidUntil int64 `protobuf:"varint,5,opt,name=ValidUntil,proto3" json:"valid_until"`
	// Resources is the list of resources an account resolves to
	Resources []*Resource `protobuf:"bytes,6,rep,name=Resources,proto3" json:"resources"`
	// Certificates contains the list of certificates to identify the account owner
	Certificates [][]byte `protobuf:"bytes,7,rep,name=Certificates,proto3" json:"certificates"`
	// MetadataURI contains a link to extra information regarding the account
	MetadataURI          string   `protobuf:"bytes,8,opt,name=MetadataURI,proto3" json:"metadata_uri"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb68cf137d2ded7a, []int{2}
}
func (m *Account) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Account.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return m.Size()
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *Account) GetName() *string {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Account) GetOwner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *Account) GetBroker() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Broker
	}
	return nil
}

func (m *Account) GetValidUntil() int64 {
	if m != nil {
		return m.ValidUntil
	}
	return 0
}

func (m *Account) GetResources() []*Resource {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *Account) GetCertificates() [][]byte {
	if m != nil {
		return m.Certificates
	}
	return nil
}

func (m *Account) GetMetadataURI() string {
	if m != nil {
		return m.MetadataURI
	}
	return ""
}

func init() {
	proto.RegisterType((*Resource)(nil), "Resource")
	proto.RegisterType((*Domain)(nil), "Domain")
	proto.RegisterType((*Account)(nil), "Account")
}

func init() { proto.RegisterFile("x/starname/types/types.proto", fileDescriptor_cb68cf137d2ded7a) }

var fileDescriptor_cb68cf137d2ded7a = []byte{
	// 526 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x4f, 0x8b, 0x13, 0x3d,
	0x1c, 0xc7, 0x9f, 0xd9, 0xe9, 0xdf, 0xb4, 0x0f, 0x4a, 0x10, 0x19, 0xa5, 0x34, 0xa5, 0xa7, 0x2a,
	0x74, 0x06, 0xaa, 0x78, 0x70, 0x41, 0xe8, 0xe8, 0xc5, 0x83, 0xae, 0x64, 0xdd, 0x3d, 0x78, 0x59,
	0xd2, 0x99, 0xec, 0x18, 0xb6, 0x33, 0x29, 0x49, 0x66, 0xd7, 0x7d, 0x27, 0x1e, 0x3d, 0xfb, 0x4a,
	0x3c, 0xfa, 0x0a, 0xa2, 0xd4, 0x83, 0x30, 0x2f, 0x61, 0x4f, 0x92, 0x34, 0xb3, 0xad, 0x5e, 0x44,
	0xc1, 0x4b, 0x27, 0xfd, 0xe4, 0x97, 0x2f, 0xc9, 0xe7, 0x97, 0x80, 0xc1, 0xbb, 0x48, 0x2a, 0x22,
	0x0a, 0x92, 0xd3, 0x48, 0x5d, 0xae, 0xa8, 0xdc, 0xfc, 0x86, 0x2b, 0xc1, 0x15, 0xbf, 0x7b, 0x2b,
	0xe3, 0x19, 0xb7, 0xc3, 0xc8, 0x8c, 0x1c, 0x1d, 0x66, 0x9c, 0x67, 0x4b, 0x1a, 0xd9, 0x7f, 0x8b,
	0xf2, 0x34, 0xba, 0x10, 0x64, 0xb5, 0xa2, 0xc2, 0xad, 0x1a, 0x1f, 0x80, 0x0e, 0xa6, 0x92, 0x97,
	0x22, 0xa1, 0xf0, 0x0e, 0xf0, 0x8f, 0xf0, 0xf3, 0xc0, 0x1b, 0x79, 0x93, 0x6e, 0xdc, 0xae, 0x34,
	0xf2, 0x4b, 0xc1, 0xb0, 0x61, 0x70, 0xb2, 0x2d, 0x0b, 0xf6, 0xec, 0x7c, 0xbf, 0xd2, 0xa8, 0x23,
	0x1c, 0xc3, 0xd7, 0xb3, 0xe3, 0x8f, 0x7b, 0xa0, 0xf5, 0x8c, 0xe7, 0x84, 0x15, 0x70, 0x00, 0x1a,
	0x2f, 0x49, 0x4e, 0x5d, 0x60, 0xa7, 0xd2, 0xa8, 0x61, 0xf6, 0x8e, 0x2d, 0x85, 0xaf, 0x40, 0x73,
	0x9e, 0xe6, 0xac, 0xb0, 0x79, 0xfd, 0xf8, 0x71, 0xa5, 0x51, 0x93, 0x18, 0x70, 0xa5, 0xd1, 0x34,
	0x63, 0xea, 0x6d, 0xb9, 0x08, 0x13, 0x9e, 0x47, 0x09, 0x97, 0x39, 0x97, 0xee, 0x33, 0x95, 0xe9,
	0x99, 0x3b, 0xf5, 0x3c, 0x49, 0xe6, 0x69, 0x2a, 0xa8, 0x94, 0x78, 0x13, 0x04, 0x0f, 0x41, 0x2b,
	0x16, 0xfc, 0x8c, 0x8a, 0xc0, 0xb7, 0x91, 0xfb, 0x95, 0x46, 0xad, 0x85, 0x25, 0x7f, 0x9e, 0xe9,
	0xa2, 0x60, 0x04, 0xc0, 0x31, 0x59, 0xb2, 0xf4, 0xa8, 0x50, 0x6c, 0x19, 0x34, 0x46, 0xde, 0xc4,
	0x8f, 0x6f, 0x54, 0x1a, 0xf5, 0xce, 0x0d, 0x3d, 0x29, 0x0d, 0xc6, 0x3b, 0x25, 0xf0, 0x3e, 0x68,
	0xbc, 0xbe, 0x5c, 0xd1, 0xa0, 0x69, 0x4f, 0x7d, 0xdb, 0x9c, 0xda, 0x84, 0x5f, 0x69, 0x04, 0x36,
	0x5e, 0xcc, 0x2c, 0xb6, 0x35, 0xe3, 0xef, 0x3e, 0x68, 0xcf, 0x93, 0x84, 0x97, 0x85, 0x82, 0xe3,
	0xda, 0x9b, 0xf3, 0x05, 0xcc, 0xee, 0x53, 0x4b, 0x70, 0x6d, 0xf4, 0x89, 0x33, 0x6a, 0x94, 0xf5,
	0x66, 0x83, 0x70, 0xd3, 0xdc, 0xb0, 0x6e, 0x6e, 0x78, 0xa8, 0x04, 0x2b, 0xb2, 0x63, 0xb2, 0x2c,
	0xa9, 0x6d, 0x90, 0xf5, 0xfd, 0xe1, 0x0b, 0xf2, 0xb6, 0xce, 0x0f, 0x2e, 0x8a, 0x6b, 0x41, 0xd6,
	0x39, 0x37, 0xe0, 0x2f, 0x9c, 0xdb, 0xa0, 0x1d, 0xe7, 0x8d, 0x7f, 0xe5, 0xbc, 0xf9, 0x7b, 0xe7,
	0x8f, 0x40, 0xb7, 0xbe, 0x80, 0x32, 0x68, 0x8d, 0xfc, 0x49, 0x6f, 0xd6, 0x0d, 0x6b, 0x12, 0xff,
	0x5f, 0x69, 0xd4, 0xad, 0xaf, 0xaa, 0xc4, 0xdb, 0x52, 0xf8, 0x10, 0xf4, 0x9f, 0x52, 0xa1, 0xd8,
	0x29, 0x4b, 0x88, 0xa2, 0x32, 0x68, 0x8f, 0xfc, 0x49, 0x3f, 0xbe, 0x59, 0x69, 0xd4, 0x4f, 0x76,
	0x38, 0xfe, 0xa9, 0x0a, 0xce, 0x40, 0xef, 0x05, 0x55, 0x24, 0x25, 0x8a, 0x98, 0xf7, 0xd2, 0xb1,
	0xed, 0xb2, 0x8b, 0x72, 0x87, 0x4f, 0xcc, 0xc3, 0xd9, 0x2d, 0x8a, 0xf7, 0x3f, 0xad, 0x87, 0xde,
	0xe7, 0xf5, 0xd0, 0xfb, 0xba, 0x1e, 0x7a, 0xef, 0xbf, 0x0d, 0xff, 0x7b, 0x73, 0x6f, 0xc7, 0x0d,
	0xe3, 0xe7, 0x53, 0x5e, 0x50, 0xf3, 0x2d, 0x64, 0xf4, 0xeb, 0x33, 0x5f, 0xb4, 0x6c, 0x83, 0x1f,
	0xfc, 0x08, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x4f, 0x55, 0xd0, 0x01, 0x04, 0x00, 0x00,
}

func (m *Resource) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Resource) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Resource) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Resource) > 0 {
		i -= len(m.Resource)
		copy(dAtA[i:], m.Resource)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Resource)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.URI) > 0 {
		i -= len(m.URI)
		copy(dAtA[i:], m.URI)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.URI)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Domain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Domain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Domain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0x2a
	}
	if m.ValidUntil != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.ValidUntil))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Broker) > 0 {
		i -= len(m.Broker)
		copy(dAtA[i:], m.Broker)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Broker)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Admin) > 0 {
		i -= len(m.Admin)
		copy(dAtA[i:], m.Admin)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Admin)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Account) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Account) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Account) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.MetadataURI) > 0 {
		i -= len(m.MetadataURI)
		copy(dAtA[i:], m.MetadataURI)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.MetadataURI)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.Certificates) > 0 {
		for iNdEx := len(m.Certificates) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Certificates[iNdEx])
			copy(dAtA[i:], m.Certificates[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Certificates[iNdEx])))
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Resources) > 0 {
		for iNdEx := len(m.Resources) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Resources[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.ValidUntil != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.ValidUntil))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Broker) > 0 {
		i -= len(m.Broker)
		copy(dAtA[i:], m.Broker)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Broker)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Name != nil {
		n1, err1 := github_com_gogo_protobuf_types.StdStringMarshalTo(*m.Name, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdString(*m.Name):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintTypes(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Domain) > 0 {
		i -= len(m.Domain)
		copy(dAtA[i:], m.Domain)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Domain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Resource) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.URI)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Resource)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Domain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Admin)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Broker)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.ValidUntil != 0 {
		n += 1 + sovTypes(uint64(m.ValidUntil))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Account) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Domain)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Name != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdString(*m.Name)
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Broker)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.ValidUntil != 0 {
		n += 1 + sovTypes(uint64(m.ValidUntil))
	}
	if len(m.Resources) > 0 {
		for _, e := range m.Resources {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.Certificates) > 0 {
		for _, b := range m.Certificates {
			l = len(b)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	l = len(m.MetadataURI)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Resource) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Resource: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Resource: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field URI", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.URI = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Resource", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Resource = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Domain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Domain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Domain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Admin", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Admin = append(m.Admin[:0], dAtA[iNdEx:postIndex]...)
			if m.Admin == nil {
				m.Admin = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Broker", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Broker = append(m.Broker[:0], dAtA[iNdEx:postIndex]...)
			if m.Broker == nil {
				m.Broker = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidUntil", wireType)
			}
			m.ValidUntil = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidUntil |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = DomainType(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Account) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Account: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Account: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Domain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Domain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Name == nil {
				m.Name = new(string)
			}
			if err := github_com_gogo_protobuf_types.StdStringUnmarshal(m.Name, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = append(m.Owner[:0], dAtA[iNdEx:postIndex]...)
			if m.Owner == nil {
				m.Owner = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Broker", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Broker = append(m.Broker[:0], dAtA[iNdEx:postIndex]...)
			if m.Broker == nil {
				m.Broker = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidUntil", wireType)
			}
			m.ValidUntil = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidUntil |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Resources", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Resources = append(m.Resources, &Resource{})
			if err := m.Resources[len(m.Resources)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Certificates", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Certificates = append(m.Certificates, make([]byte, postIndex-iNdEx))
			copy(m.Certificates[len(m.Certificates)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetadataURI", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MetadataURI = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)