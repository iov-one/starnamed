// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: x/starname/types/types.proto

package types

import (
	bytes "bytes"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Resource defines a resource owned by an account
type Resource struct {
	// URI defines the ID of the resource
	Uri string `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty" yaml:"uri"`
	// Resource is the resource
	Resource string `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty" yaml:"resource"`
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

func (m *Resource) GetUri() string {
	if m != nil {
		return m.Uri
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
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" yaml:"name"`
	// Admin is the owner of the domain
	Admin  github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,2,opt,name=admin,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"admin,omitempty" yaml:"admin"`
	Broker github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=broker,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"broker,omitempty" yaml:"broker"`
	// ValidUntil is a unix timestamp defines the time when the domain will become invalid in seconds
	ValidUntil int64 `protobuf:"varint,4,opt,name=valid_until,json=validUntil,proto3" json:"valid_until,omitempty" yaml:"valid_until"`
	// Type defines the type of the domain
	Type DomainType `protobuf:"bytes,5,opt,name=type,proto3,casttype=DomainType" json:"type,omitempty" yaml:"type"`
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
	Domain string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty" yaml:"domain"`
	// Name is the name of the account
	Name *string `protobuf:"bytes,2,opt,name=name,proto3,wktptr" json:"name,omitempty" yaml:"name"`
	// Owner is the address that owns the account
	Owner github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=owner,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"owner,omitempty" yaml:"owner"`
	// Broker identifies an entity that facilitated the transaction of the account and can be empty
	Broker github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,4,opt,name=broker,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"broker,omitempty" yaml:"broker"`
	// ValidUntil defines a unix timestamp of the expiration of the account in seconds
	ValidUntil int64 `protobuf:"varint,5,opt,name=valid_until,json=validUntil,proto3" json:"valid_until,omitempty" yaml:"valid_until"`
	// Resources is the list of resources an account resolves to
	Resources []*Resource `protobuf:"bytes,6,rep,name=resources,proto3" json:"resources,omitempty" yaml:"resources"`
	// Certificates contains the list of certificates to identify the account owner
	Certificates [][]byte `protobuf:"bytes,7,rep,name=certificates,proto3" json:"certificates,omitempty" yaml:"certificates"`
	// MetadataURI contains a link to extra information regarding the account
	ValidURI string `protobuf:"bytes,8,opt,name=metadata_uri,json=metadataUri,proto3" json:"metadata_uri,omitempty" yaml:"metadata_uri"`
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

func (m *Account) GetValidURI() string {
	if m != nil {
		return m.ValidURI
	}
	return ""
}

func init() {
	proto.RegisterType((*Resource)(nil), "starnamed.x.starname.v1beta1.Resource")
	proto.RegisterType((*Domain)(nil), "starnamed.x.starname.v1beta1.Domain")
	proto.RegisterType((*Account)(nil), "starnamed.x.starname.v1beta1.Account")
}

func init() { proto.RegisterFile("x/starname/types/types.proto", fileDescriptor_cb68cf137d2ded7a) }

var fileDescriptor_cb68cf137d2ded7a = []byte{
	// 595 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0x5e, 0xd6, 0xae, 0xeb, 0xdc, 0xc2, 0xc0, 0x9b, 0x20, 0x9a, 0x46, 0x52, 0x79, 0x12, 0x2a,
	0x87, 0x25, 0xda, 0x38, 0x20, 0xc1, 0xa9, 0x11, 0x20, 0x21, 0x6e, 0x86, 0x0d, 0x69, 0x12, 0x9a,
	0xdc, 0xc4, 0x0b, 0xd6, 0x9a, 0xb8, 0x72, 0x9c, 0x6e, 0xfd, 0x17, 0xfc, 0x04, 0x8e, 0xfc, 0x14,
	0xc4, 0x01, 0xed, 0xc8, 0x29, 0xa0, 0xf6, 0x1f, 0xe4, 0xc8, 0x09, 0xc5, 0x4e, 0xda, 0x8c, 0x03,
	0x42, 0x93, 0xb8, 0xb4, 0x7e, 0xef, 0x7d, 0xef, 0xab, 0xfb, 0x7d, 0xef, 0x19, 0xec, 0x5e, 0xba,
	0x89, 0x24, 0x22, 0x26, 0x11, 0x75, 0xe5, 0x74, 0x4c, 0x13, 0xfd, 0xe9, 0x8c, 0x05, 0x97, 0x1c,
	0xee, 0x56, 0xb5, 0xc0, 0xb9, 0x74, 0xaa, 0xb3, 0x33, 0x39, 0x18, 0x52, 0x49, 0x0e, 0x76, 0xb6,
	0x43, 0x1e, 0x72, 0x05, 0x74, 0x8b, 0x93, 0xee, 0xd9, 0xb1, 0x42, 0xce, 0xc3, 0x11, 0x75, 0x55,
	0x34, 0x4c, 0xcf, 0xdc, 0x0b, 0x41, 0xc6, 0x63, 0x2a, 0x4a, 0x4e, 0xf4, 0x1e, 0xb4, 0x31, 0x4d,
	0x78, 0x2a, 0x7c, 0x0a, 0x7b, 0xa0, 0x91, 0x0a, 0x66, 0x1a, 0x3d, 0xa3, 0xbf, 0xe1, 0xdd, 0xce,
	0x33, 0x1b, 0x4c, 0x49, 0x34, 0x7a, 0x8a, 0x52, 0xc1, 0x10, 0x2e, 0x4a, 0xd0, 0x05, 0x6d, 0x51,
	0xa2, 0xcd, 0x55, 0x05, 0xdb, 0xca, 0x33, 0x7b, 0x53, 0xc3, 0xaa, 0x0a, 0xc2, 0x0b, 0x10, 0xfa,
	0xb6, 0x0a, 0x5a, 0xcf, 0x79, 0x44, 0x58, 0x0c, 0xf7, 0x40, 0xb3, 0xb8, 0x6f, 0x49, 0xbf, 0x99,
	0x67, 0x76, 0x47, 0xf7, 0x15, 0x59, 0x84, 0x55, 0x11, 0xbe, 0x03, 0x6b, 0x24, 0x88, 0x58, 0xac,
	0xd8, 0xbb, 0xde, 0x20, 0xcf, 0xec, 0xae, 0x46, 0xa9, 0x34, 0xfa, 0x95, 0xd9, 0xfb, 0x21, 0x93,
	0x1f, 0xd2, 0xa1, 0xe3, 0xf3, 0xc8, 0xf5, 0x79, 0x12, 0xf1, 0xa4, 0xfc, 0xda, 0x4f, 0x82, 0xf3,
	0x52, 0xaf, 0x81, 0xef, 0x0f, 0x82, 0x40, 0xd0, 0x24, 0xc1, 0x9a, 0x0f, 0x9e, 0x80, 0xd6, 0x50,
	0xf0, 0x73, 0x2a, 0xcc, 0x86, 0x62, 0xf6, 0xf2, 0xcc, 0xbe, 0xa5, 0x99, 0x75, 0xfe, 0x06, 0xd4,
	0x25, 0x23, 0x7c, 0x02, 0x3a, 0x13, 0x32, 0x62, 0xc1, 0x69, 0x1a, 0x4b, 0x36, 0x32, 0x9b, 0x3d,
	0xa3, 0xdf, 0xf0, 0xee, 0xe5, 0x99, 0x0d, 0xf5, 0x0f, 0xd4, 0x8a, 0x08, 0x03, 0x15, 0x1d, 0x15,
	0x01, 0x3c, 0x00, 0xcd, 0x82, 0xd4, 0x5c, 0x53, 0x92, 0x3c, 0x58, 0x4a, 0x52, 0x64, 0x8b, 0x0b,
	0x01, 0xad, 0xdd, 0xdb, 0xe9, 0x98, 0x62, 0x05, 0x45, 0x5f, 0x9b, 0x60, 0x7d, 0xe0, 0xfb, 0x3c,
	0x8d, 0x25, 0x7c, 0x04, 0x5a, 0x81, 0xaa, 0x97, 0x9a, 0xde, 0x5d, 0xfe, 0x27, 0x9d, 0x47, 0xb8,
	0x04, 0xc0, 0x17, 0xa5, 0xf8, 0x85, 0xac, 0x9d, 0xc3, 0x5d, 0x47, 0x4f, 0x85, 0x53, 0x4d, 0x85,
	0xf3, 0x46, 0x0a, 0x16, 0x87, 0xc7, 0x64, 0x94, 0x52, 0x65, 0x69, 0xdd, 0x9a, 0x4f, 0x3f, 0x6c,
	0x63, 0x69, 0x0f, 0xbf, 0x88, 0x17, 0x22, 0xd6, 0xec, 0x51, 0xe9, 0x9b, 0xd8, 0xa3, 0x1a, 0x6b,
	0xf6, 0x34, 0xff, 0xb7, 0x3d, 0x6b, 0xff, 0x6c, 0xcf, 0x09, 0xd8, 0xa8, 0x06, 0x39, 0x31, 0x5b,
	0xbd, 0x46, 0xbf, 0x73, 0xf8, 0xd0, 0xf9, 0xdb, 0x0e, 0x3a, 0xd5, 0x2a, 0x79, 0xdb, 0x79, 0x66,
	0xdf, 0xb9, 0xbe, 0x16, 0x09, 0xc2, 0x4b, 0x3a, 0xf8, 0x0c, 0x74, 0x7d, 0x2a, 0x24, 0x3b, 0x63,
	0x3e, 0x91, 0x34, 0x31, 0xd7, 0x7b, 0x8d, 0x7e, 0xd7, 0xbb, 0x9f, 0x67, 0xf6, 0x96, 0x6e, 0xab,
	0x57, 0x11, 0xbe, 0x06, 0x86, 0x2f, 0x41, 0x37, 0xa2, 0x92, 0x04, 0x44, 0x92, 0xd3, 0x62, 0x63,
	0xdb, 0xca, 0xfe, 0xbd, 0x59, 0x66, 0xb7, 0x8f, 0xd5, 0xf5, 0xf1, 0xab, 0x25, 0x51, 0x1d, 0x89,
	0x70, 0xa7, 0x0a, 0x8f, 0x04, 0xf3, 0x5e, 0x7f, 0x9e, 0x59, 0xc6, 0x97, 0x99, 0x65, 0x5c, 0xcd,
	0x2c, 0xe3, 0xe7, 0xcc, 0x32, 0x3e, 0xce, 0xad, 0x95, 0xab, 0xb9, 0xb5, 0xf2, 0x7d, 0x6e, 0xad,
	0x9c, 0xd4, 0xe5, 0x66, 0x7c, 0xb2, 0xcf, 0x63, 0xba, 0x78, 0x9d, 0x02, 0xf7, 0xcf, 0x97, 0x6a,
	0xd8, 0x52, 0xc3, 0xf4, 0xf8, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x73, 0x85, 0x8c, 0xd9, 0xc4,
	0x04, 0x00, 0x00,
}

func (this *Resource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Resource)
	if !ok {
		that2, ok := that.(Resource)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Uri != that1.Uri {
		return false
	}
	if this.Resource != that1.Resource {
		return false
	}
	return true
}
func (this *Domain) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Domain)
	if !ok {
		that2, ok := that.(Domain)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if !bytes.Equal(this.Admin, that1.Admin) {
		return false
	}
	if !bytes.Equal(this.Broker, that1.Broker) {
		return false
	}
	if this.ValidUntil != that1.ValidUntil {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	return true
}
func (this *Account) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Account)
	if !ok {
		that2, ok := that.(Account)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Domain != that1.Domain {
		return false
	}
	if this.Name != nil && that1.Name != nil {
		if *this.Name != *that1.Name {
			return false
		}
	} else if this.Name != nil {
		return false
	} else if that1.Name != nil {
		return false
	}
	if !bytes.Equal(this.Owner, that1.Owner) {
		return false
	}
	if !bytes.Equal(this.Broker, that1.Broker) {
		return false
	}
	if this.ValidUntil != that1.ValidUntil {
		return false
	}
	if len(this.Resources) != len(that1.Resources) {
		return false
	}
	for i := range this.Resources {
		if !this.Resources[i].Equal(that1.Resources[i]) {
			return false
		}
	}
	if len(this.Certificates) != len(that1.Certificates) {
		return false
	}
	for i := range this.Certificates {
		if !bytes.Equal(this.Certificates[i], that1.Certificates[i]) {
			return false
		}
	}
	if this.ValidURI != that1.ValidURI {
		return false
	}
	return true
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
	if len(m.Resource) > 0 {
		i -= len(m.Resource)
		copy(dAtA[i:], m.Resource)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Resource)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Uri) > 0 {
		i -= len(m.Uri)
		copy(dAtA[i:], m.Uri)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Uri)))
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
	if len(m.ValidURI) > 0 {
		i -= len(m.ValidURI)
		copy(dAtA[i:], m.ValidURI)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidURI)))
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
	l = len(m.Uri)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Resource)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
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
	l = len(m.ValidURI)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Uri", wireType)
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
			m.Uri = string(dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field ValidURI", wireType)
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
			m.ValidURI = string(dAtA[iNdEx:postIndex])
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
