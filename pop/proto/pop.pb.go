// Code generated by protoc-gen-go.
// source: pop.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	pop.proto

It has these top-level messages:
	Container
	ContainerConfig
	ContainerList
	Credentials
	Endpoint
	Filter
	Flavour
	FlavourList
	Image
	ImageList
	Infos
	Ip
	Metadata
	Network
	NetworkList
	NewMetadata
	Subnet
	Token
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Container_Status int32

const (
	Container_UNAVAILABLE Container_Status = 0
	Container_CREATED     Container_Status = 1
	Container_RUNNING     Container_Status = 2
	Container_EXITED      Container_Status = 3
	Container_FAILED      Container_Status = 4
)

var Container_Status_name = map[int32]string{
	0: "UNAVAILABLE",
	1: "CREATED",
	2: "RUNNING",
	3: "EXITED",
	4: "FAILED",
}
var Container_Status_value = map[string]int32{
	"UNAVAILABLE": 0,
	"CREATED":     1,
	"RUNNING":     2,
	"EXITED":      3,
	"FAILED":      4,
}

func (x Container_Status) String() string {
	return proto1.EnumName(Container_Status_name, int32(x))
}
func (Container_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Container struct {
	Id             string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Names          []string             `protobuf:"bytes,2,rep,name=names" json:"names,omitempty"`
	ImageId        string               `protobuf:"bytes,3,opt,name=image_id,json=imageId" json:"image_id,omitempty"`
	FlavourId      string               `protobuf:"bytes,4,opt,name=flavour_id,json=flavourId" json:"flavour_id,omitempty"`
	Command        string               `protobuf:"bytes,5,opt,name=command" json:"command,omitempty"`
	Created        int64                `protobuf:"varint,6,opt,name=created" json:"created,omitempty"`
	Status         Container_Status     `protobuf:"varint,7,opt,name=status,enum=vim_pop.Container_Status" json:"status,omitempty"`
	ExtendedStatus string               `protobuf:"bytes,8,opt,name=extended_status,json=extendedStatus" json:"extended_status,omitempty"`
	Endpoints      map[string]*Endpoint `protobuf:"bytes,9,rep,name=endpoints" json:"endpoints,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Md             *Metadata            `protobuf:"bytes,10,opt,name=md" json:"md,omitempty"`
}

func (m *Container) Reset()                    { *m = Container{} }
func (m *Container) String() string            { return proto1.CompactTextString(m) }
func (*Container) ProtoMessage()               {}
func (*Container) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Container) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Container) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
}

func (m *Container) GetImageId() string {
	if m != nil {
		return m.ImageId
	}
	return ""
}

func (m *Container) GetFlavourId() string {
	if m != nil {
		return m.FlavourId
	}
	return ""
}

func (m *Container) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *Container) GetCreated() int64 {
	if m != nil {
		return m.Created
	}
	return 0
}

func (m *Container) GetStatus() Container_Status {
	if m != nil {
		return m.Status
	}
	return Container_UNAVAILABLE
}

func (m *Container) GetExtendedStatus() string {
	if m != nil {
		return m.ExtendedStatus
	}
	return ""
}

func (m *Container) GetEndpoints() map[string]*Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

func (m *Container) GetMd() *Metadata {
	if m != nil {
		return m.Md
	}
	return nil
}

type ContainerConfig struct {
	Name      string               `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ImageId   string               `protobuf:"bytes,2,opt,name=image_id,json=imageId" json:"image_id,omitempty"`
	FlavourId string               `protobuf:"bytes,3,opt,name=flavour_id,json=flavourId" json:"flavour_id,omitempty"`
	Endpoints map[string]*Endpoint `protobuf:"bytes,4,rep,name=endpoints" json:"endpoints,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ContainerConfig) Reset()                    { *m = ContainerConfig{} }
func (m *ContainerConfig) String() string            { return proto1.CompactTextString(m) }
func (*ContainerConfig) ProtoMessage()               {}
func (*ContainerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ContainerConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ContainerConfig) GetImageId() string {
	if m != nil {
		return m.ImageId
	}
	return ""
}

func (m *ContainerConfig) GetFlavourId() string {
	if m != nil {
		return m.FlavourId
	}
	return ""
}

func (m *ContainerConfig) GetEndpoints() map[string]*Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

type ContainerList struct {
	List []*Container `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *ContainerList) Reset()                    { *m = ContainerList{} }
func (m *ContainerList) String() string            { return proto1.CompactTextString(m) }
func (*ContainerList) ProtoMessage()               {}
func (*ContainerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ContainerList) GetList() []*Container {
	if m != nil {
		return m.List
	}
	return nil
}

// Credentials represents the login credentials for a given user.
type Credentials struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Credentials) Reset()                    { *m = Credentials{} }
func (m *Credentials) String() string            { return proto1.CompactTextString(m) }
func (*Credentials) ProtoMessage()               {}
func (*Credentials) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Credentials) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Credentials) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Endpoint struct {
	NetId      string `protobuf:"bytes,1,opt,name=net_id,json=netId" json:"net_id,omitempty"`
	EndpointId string `protobuf:"bytes,2,opt,name=endpoint_id,json=endpointId" json:"endpoint_id,omitempty"`
	Ipv4       *Ip    `protobuf:"bytes,3,opt,name=ipv4" json:"ipv4,omitempty"`
	Ipv6       *Ip    `protobuf:"bytes,4,opt,name=ipv6" json:"ipv6,omitempty"`
	Mac        string `protobuf:"bytes,5,opt,name=mac" json:"mac,omitempty"`
}

func (m *Endpoint) Reset()                    { *m = Endpoint{} }
func (m *Endpoint) String() string            { return proto1.CompactTextString(m) }
func (*Endpoint) ProtoMessage()               {}
func (*Endpoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Endpoint) GetNetId() string {
	if m != nil {
		return m.NetId
	}
	return ""
}

func (m *Endpoint) GetEndpointId() string {
	if m != nil {
		return m.EndpointId
	}
	return ""
}

func (m *Endpoint) GetIpv4() *Ip {
	if m != nil {
		return m.Ipv4
	}
	return nil
}

func (m *Endpoint) GetIpv6() *Ip {
	if m != nil {
		return m.Ipv6
	}
	return nil
}

func (m *Endpoint) GetMac() string {
	if m != nil {
		return m.Mac
	}
	return ""
}

type Filter struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto1.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Filter) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Flavour struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *Flavour) Reset()                    { *m = Flavour{} }
func (m *Flavour) String() string            { return proto1.CompactTextString(m) }
func (*Flavour) ProtoMessage()               {}
func (*Flavour) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Flavour) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Flavour) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type FlavourList struct {
	List []*Flavour `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *FlavourList) Reset()                    { *m = FlavourList{} }
func (m *FlavourList) String() string            { return proto1.CompactTextString(m) }
func (*FlavourList) ProtoMessage()               {}
func (*FlavourList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *FlavourList) GetList() []*Flavour {
	if m != nil {
		return m.List
	}
	return nil
}

type Image struct {
	Id      string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Names   []string `protobuf:"bytes,2,rep,name=names" json:"names,omitempty"`
	Created int64    `protobuf:"varint,3,opt,name=created" json:"created,omitempty"`
}

func (m *Image) Reset()                    { *m = Image{} }
func (m *Image) String() string            { return proto1.CompactTextString(m) }
func (*Image) ProtoMessage()               {}
func (*Image) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Image) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Image) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
}

func (m *Image) GetCreated() int64 {
	if m != nil {
		return m.Created
	}
	return 0
}

type ImageList struct {
	List []*Image `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *ImageList) Reset()                    { *m = ImageList{} }
func (m *ImageList) String() string            { return proto1.CompactTextString(m) }
func (*ImageList) ProtoMessage()               {}
func (*ImageList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ImageList) GetList() []*Image {
	if m != nil {
		return m.List
	}
	return nil
}

type Infos struct {
	Type      string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Infos) Reset()                    { *m = Infos{} }
func (m *Infos) String() string            { return proto1.CompactTextString(m) }
func (*Infos) ProtoMessage()               {}
func (*Infos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Infos) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Infos) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Infos) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type Ip struct {
	Address string  `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	Subnet  *Subnet `protobuf:"bytes,2,opt,name=subnet" json:"subnet,omitempty"`
}

func (m *Ip) Reset()                    { *m = Ip{} }
func (m *Ip) String() string            { return proto1.CompactTextString(m) }
func (*Ip) ProtoMessage()               {}
func (*Ip) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Ip) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Ip) GetSubnet() *Subnet {
	if m != nil {
		return m.Subnet
	}
	return nil
}

// Metadata contains a key-value set of metadata
// pairs, that will be exposed to the underlying container.
type Metadata struct {
	Entries map[string]string `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Metadata) Reset()                    { *m = Metadata{} }
func (m *Metadata) String() string            { return proto1.CompactTextString(m) }
func (*Metadata) ProtoMessage()               {}
func (*Metadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Metadata) GetEntries() map[string]string {
	if m != nil {
		return m.Entries
	}
	return nil
}

type Network struct {
	Id       string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name     string    `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	External bool      `protobuf:"varint,3,opt,name=external" json:"external,omitempty"`
	Subnets  []*Subnet `protobuf:"bytes,4,rep,name=subnets" json:"subnets,omitempty"`
}

func (m *Network) Reset()                    { *m = Network{} }
func (m *Network) String() string            { return proto1.CompactTextString(m) }
func (*Network) ProtoMessage()               {}
func (*Network) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Network) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Network) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Network) GetExternal() bool {
	if m != nil {
		return m.External
	}
	return false
}

func (m *Network) GetSubnets() []*Subnet {
	if m != nil {
		return m.Subnets
	}
	return nil
}

type NetworkList struct {
	List []*Network `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"`
}

func (m *NetworkList) Reset()                    { *m = NetworkList{} }
func (m *NetworkList) String() string            { return proto1.CompactTextString(m) }
func (*NetworkList) ProtoMessage()               {}
func (*NetworkList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *NetworkList) GetList() []*Network {
	if m != nil {
		return m.List
	}
	return nil
}

type NewMetadata struct {
	Id string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Md *Metadata `protobuf:"bytes,2,opt,name=md" json:"md,omitempty"`
}

func (m *NewMetadata) Reset()                    { *m = NewMetadata{} }
func (m *NewMetadata) String() string            { return proto1.CompactTextString(m) }
func (*NewMetadata) ProtoMessage()               {}
func (*NewMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *NewMetadata) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NewMetadata) GetMd() *Metadata {
	if m != nil {
		return m.Md
	}
	return nil
}

type Subnet struct {
	Cidr    string `protobuf:"bytes,1,opt,name=cidr" json:"cidr,omitempty"`
	Gateway string `protobuf:"bytes,2,opt,name=gateway" json:"gateway,omitempty"`
}

func (m *Subnet) Reset()                    { *m = Subnet{} }
func (m *Subnet) String() string            { return proto1.CompactTextString(m) }
func (*Subnet) ProtoMessage()               {}
func (*Subnet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *Subnet) GetCidr() string {
	if m != nil {
		return m.Cidr
	}
	return ""
}

func (m *Subnet) GetGateway() string {
	if m != nil {
		return m.Gateway
	}
	return ""
}

// Token is a token generated by the server after a successful login.
// This token should be set as metadata, to authenticate every other
type Token struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto1.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *Token) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto1.RegisterType((*Container)(nil), "vim_pop.Container")
	proto1.RegisterType((*ContainerConfig)(nil), "vim_pop.ContainerConfig")
	proto1.RegisterType((*ContainerList)(nil), "vim_pop.ContainerList")
	proto1.RegisterType((*Credentials)(nil), "vim_pop.Credentials")
	proto1.RegisterType((*Endpoint)(nil), "vim_pop.Endpoint")
	proto1.RegisterType((*Filter)(nil), "vim_pop.Filter")
	proto1.RegisterType((*Flavour)(nil), "vim_pop.Flavour")
	proto1.RegisterType((*FlavourList)(nil), "vim_pop.FlavourList")
	proto1.RegisterType((*Image)(nil), "vim_pop.Image")
	proto1.RegisterType((*ImageList)(nil), "vim_pop.ImageList")
	proto1.RegisterType((*Infos)(nil), "vim_pop.Infos")
	proto1.RegisterType((*Ip)(nil), "vim_pop.Ip")
	proto1.RegisterType((*Metadata)(nil), "vim_pop.Metadata")
	proto1.RegisterType((*Network)(nil), "vim_pop.Network")
	proto1.RegisterType((*NetworkList)(nil), "vim_pop.NetworkList")
	proto1.RegisterType((*NewMetadata)(nil), "vim_pop.NewMetadata")
	proto1.RegisterType((*Subnet)(nil), "vim_pop.Subnet")
	proto1.RegisterType((*Token)(nil), "vim_pop.Token")
	proto1.RegisterEnum("vim_pop.Container_Status", Container_Status_name, Container_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Pop service

type PopClient interface {
	// Containers returns the containers available in the PoP, either
	// created or running.
	Containers(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ContainerList, error)
	// Flavours returns the available flavours.
	// This doesn't make much sense with containers, but it's here to
	// better abstract the PoP.
	Flavours(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*FlavourList, error)
	// Images returns the images available in the PoP.
	Images(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ImageList, error)
	// Networks returns the available retworks in the PoP.
	Networks(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*NetworkList, error)
	// Create creates a new container as described.
	Create(ctx context.Context, in *ContainerConfig, opts ...grpc.CallOption) (*Container, error)
	// Delete stops and deletes the container identified by the given filter.
	Delete(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	// Metadata adds the given metadata values to the container that matches with the ID.
	// An empty value for a key means that the key will be removed from the metadata.
	Metadata(ctx context.Context, in *NewMetadata, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	// Start starts the container identified by the given filter.
	// Any metadata key stored in the server will be passed to the newly instantiated container.
	Start(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Container, error)
	// Stop starts the container identified by the given filter.
	Stop(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	// Login logs an user in and sets up a session.
	// The returned token should be set into the metadata
	// of the gRPC session with key "token" to authenticate your client.
	Login(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error)
	// Logout invalids the current token.
	Logout(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	// Info can be used to check if the Pop is alive and if your credentials to this service are valid.
	// It also returns informations about this server.
	Info(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*Infos, error)
}

type popClient struct {
	cc *grpc.ClientConn
}

func NewPopClient(cc *grpc.ClientConn) PopClient {
	return &popClient{cc}
}

func (c *popClient) Containers(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ContainerList, error) {
	out := new(ContainerList)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Containers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Flavours(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*FlavourList, error) {
	out := new(FlavourList)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Flavours", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Images(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ImageList, error) {
	out := new(ImageList)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Images", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Networks(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*NetworkList, error) {
	out := new(NetworkList)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Networks", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Create(ctx context.Context, in *ContainerConfig, opts ...grpc.CallOption) (*Container, error) {
	out := new(Container)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Delete(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Metadata(ctx context.Context, in *NewMetadata, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Metadata", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Start(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Container, error) {
	out := new(Container)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Stop(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Login(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Login", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Logout(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Logout", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Info(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*Infos, error) {
	out := new(Infos)
	err := grpc.Invoke(ctx, "/vim_pop.Pop/Info", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Pop service

type PopServer interface {
	// Containers returns the containers available in the PoP, either
	// created or running.
	Containers(context.Context, *Filter) (*ContainerList, error)
	// Flavours returns the available flavours.
	// This doesn't make much sense with containers, but it's here to
	// better abstract the PoP.
	Flavours(context.Context, *Filter) (*FlavourList, error)
	// Images returns the images available in the PoP.
	Images(context.Context, *Filter) (*ImageList, error)
	// Networks returns the available retworks in the PoP.
	Networks(context.Context, *Filter) (*NetworkList, error)
	// Create creates a new container as described.
	Create(context.Context, *ContainerConfig) (*Container, error)
	// Delete stops and deletes the container identified by the given filter.
	Delete(context.Context, *Filter) (*google_protobuf.Empty, error)
	// Metadata adds the given metadata values to the container that matches with the ID.
	// An empty value for a key means that the key will be removed from the metadata.
	Metadata(context.Context, *NewMetadata) (*google_protobuf.Empty, error)
	// Start starts the container identified by the given filter.
	// Any metadata key stored in the server will be passed to the newly instantiated container.
	Start(context.Context, *Filter) (*Container, error)
	// Stop starts the container identified by the given filter.
	Stop(context.Context, *Filter) (*google_protobuf.Empty, error)
	// Login logs an user in and sets up a session.
	// The returned token should be set into the metadata
	// of the gRPC session with key "token" to authenticate your client.
	Login(context.Context, *Credentials) (*Token, error)
	// Logout invalids the current token.
	Logout(context.Context, *google_protobuf.Empty) (*google_protobuf.Empty, error)
	// Info can be used to check if the Pop is alive and if your credentials to this service are valid.
	// It also returns informations about this server.
	Info(context.Context, *google_protobuf.Empty) (*Infos, error)
}

func RegisterPopServer(s *grpc.Server, srv PopServer) {
	s.RegisterService(&_Pop_serviceDesc, srv)
}

func _Pop_Containers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Containers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Containers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Containers(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Flavours_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Flavours(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Flavours",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Flavours(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Images_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Images(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Images",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Images(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Networks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Networks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Networks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Networks(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContainerConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Create(ctx, req.(*ContainerConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Delete(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Metadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewMetadata)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Metadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Metadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Metadata(ctx, req.(*NewMetadata))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Start(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Stop(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Login(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Logout(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pop_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PopServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vim_pop.Pop/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Info(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pop_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vim_pop.Pop",
	HandlerType: (*PopServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Containers",
			Handler:    _Pop_Containers_Handler,
		},
		{
			MethodName: "Flavours",
			Handler:    _Pop_Flavours_Handler,
		},
		{
			MethodName: "Images",
			Handler:    _Pop_Images_Handler,
		},
		{
			MethodName: "Networks",
			Handler:    _Pop_Networks_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Pop_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Pop_Delete_Handler,
		},
		{
			MethodName: "Metadata",
			Handler:    _Pop_Metadata_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _Pop_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Pop_Stop_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Pop_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Pop_Logout_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Pop_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pop.proto",
}

func init() { proto1.RegisterFile("pop.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 990 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x56, 0x5f, 0x6f, 0xe3, 0x44,
	0x10, 0xc7, 0x76, 0x6c, 0xc7, 0x63, 0x48, 0xc3, 0xaa, 0x20, 0x5f, 0xc4, 0x41, 0xce, 0x42, 0x34,
	0x48, 0x34, 0xa7, 0xcb, 0x41, 0xa9, 0xee, 0x05, 0x72, 0xa9, 0x5b, 0x59, 0xe4, 0x02, 0x72, 0xef,
	0x10, 0xe2, 0xa5, 0xda, 0xab, 0xb7, 0x91, 0xd5, 0xd8, 0x6b, 0xd9, 0x9b, 0x96, 0xbc, 0xf1, 0x29,
	0xf8, 0x98, 0xbc, 0xf2, 0x84, 0x84, 0x76, 0xbd, 0xfe, 0x93, 0xe6, 0x8f, 0xfa, 0xc0, 0x53, 0x76,
	0x76, 0x66, 0x76, 0x66, 0x7e, 0x33, 0xf3, 0x8b, 0xc1, 0x4a, 0x69, 0x3a, 0x4c, 0x33, 0xca, 0x28,
	0x32, 0xef, 0xa2, 0xf8, 0x2a, 0xa5, 0x69, 0xcf, 0x26, 0x71, 0xca, 0x56, 0xc5, 0xad, 0xfb, 0xaf,
	0x06, 0xd6, 0x84, 0x26, 0x0c, 0x47, 0x09, 0xc9, 0x50, 0x07, 0xd4, 0x28, 0x74, 0x94, 0xbe, 0x32,
	0xb0, 0x02, 0x35, 0x0a, 0xd1, 0x21, 0xe8, 0x09, 0x8e, 0x49, 0xee, 0xa8, 0x7d, 0x6d, 0x60, 0x05,
	0x85, 0x80, 0x9e, 0x40, 0x3b, 0x8a, 0xf1, 0x9c, 0x5c, 0x45, 0xa1, 0xa3, 0x09, 0x5b, 0x53, 0xc8,
	0x7e, 0x88, 0x9e, 0x02, 0xdc, 0x2c, 0xf0, 0x1d, 0x5d, 0x66, 0x5c, 0xd9, 0x12, 0x4a, 0x4b, 0xde,
	0xf8, 0x21, 0x72, 0xc0, 0xbc, 0xa6, 0x71, 0x8c, 0x93, 0xd0, 0xd1, 0x0b, 0x47, 0x29, 0x0a, 0x4d,
	0x46, 0x30, 0x23, 0xa1, 0x63, 0xf4, 0x95, 0x81, 0x16, 0x94, 0x22, 0x7a, 0x01, 0x46, 0xce, 0x30,
	0x5b, 0xe6, 0x8e, 0xd9, 0x57, 0x06, 0x9d, 0xd1, 0x93, 0xa1, 0x2c, 0x64, 0x58, 0xe5, 0x3d, 0xbc,
	0x14, 0x06, 0x81, 0x34, 0x44, 0x47, 0x70, 0x40, 0xfe, 0x60, 0x24, 0x09, 0x49, 0x78, 0x25, 0x7d,
	0xdb, 0x22, 0x5c, 0xa7, 0xbc, 0x2e, 0x1c, 0xd0, 0x0f, 0x60, 0x91, 0x24, 0x4c, 0x69, 0x94, 0xb0,
	0xdc, 0xb1, 0xfa, 0xda, 0xc0, 0x1e, 0x3d, 0xdb, 0xf2, 0xbc, 0x57, 0xda, 0x78, 0x09, 0xcb, 0x56,
	0x41, 0xed, 0x83, 0x9e, 0x81, 0x1a, 0x87, 0x0e, 0xf4, 0x95, 0x81, 0x3d, 0xfa, 0xb8, 0xf2, 0x7c,
	0x43, 0x18, 0x0e, 0x31, 0xc3, 0x81, 0x1a, 0x87, 0xbd, 0x9f, 0xa1, 0xb3, 0xee, 0x8f, 0xba, 0xa0,
	0xdd, 0x92, 0x95, 0x84, 0x99, 0x1f, 0xd1, 0x11, 0xe8, 0x77, 0x78, 0xb1, 0x24, 0x8e, 0xfa, 0xe0,
	0xa5, 0xd2, 0x33, 0x28, 0xf4, 0xaf, 0xd4, 0x53, 0xc5, 0xfd, 0x09, 0x0c, 0x99, 0xfe, 0x01, 0xd8,
	0xef, 0x66, 0xe3, 0x5f, 0xc7, 0xfe, 0x74, 0xfc, 0x7a, 0xea, 0x75, 0x3f, 0x40, 0x36, 0x98, 0x93,
	0xc0, 0x1b, 0xbf, 0xf5, 0xce, 0xba, 0x0a, 0x17, 0x82, 0x77, 0xb3, 0x99, 0x3f, 0xbb, 0xe8, 0xaa,
	0x08, 0xc0, 0xf0, 0x7e, 0xf3, 0xb9, 0x42, 0xe3, 0xe7, 0xf3, 0xb1, 0x3f, 0xf5, 0xce, 0xba, 0x2d,
	0xf7, 0x1f, 0x05, 0x0e, 0xaa, 0x42, 0x27, 0x34, 0xb9, 0x89, 0xe6, 0x08, 0x41, 0x8b, 0x37, 0x5a,
	0x26, 0x28, 0xce, 0x6b, 0x3d, 0x57, 0xf7, 0xf5, 0x5c, 0x7b, 0xd8, 0x73, 0xaf, 0x89, 0x71, 0x4b,
	0x60, 0x7c, 0xb4, 0x89, 0x71, 0x11, 0x7a, 0x37, 0xd2, 0xff, 0x3f, 0x8c, 0xdf, 0xc3, 0x47, 0x55,
	0xf4, 0x69, 0x94, 0x33, 0xf4, 0x15, 0xb4, 0x16, 0x51, 0xce, 0x1c, 0x45, 0xe4, 0x88, 0x36, 0x73,
	0x0c, 0x84, 0xde, 0xf5, 0xc0, 0x9e, 0x64, 0x24, 0x24, 0x09, 0x8b, 0xf0, 0x22, 0x47, 0x3d, 0x68,
	0x2f, 0x73, 0x92, 0x35, 0x10, 0xab, 0x64, 0xae, 0x4b, 0x71, 0x9e, 0xdf, 0xd3, 0xac, 0x44, 0xad,
	0x92, 0xdd, 0xbf, 0x14, 0x68, 0x97, 0x79, 0xa1, 0x4f, 0xc0, 0x48, 0x08, 0xbb, 0xaa, 0x96, 0x4f,
	0x4f, 0x08, 0xf3, 0x43, 0xf4, 0x05, 0xd8, 0x25, 0x02, 0x35, 0xf0, 0x50, 0x5e, 0x09, 0x83, 0x56,
	0x94, 0xde, 0x7d, 0x2b, 0x50, 0xb7, 0x47, 0x76, 0x95, 0xb3, 0x9f, 0x06, 0x42, 0x21, 0x0d, 0x4e,
	0xc4, 0x2a, 0x6e, 0x31, 0x38, 0xe1, 0x28, 0xc6, 0xf8, 0x5a, 0xae, 0x23, 0x3f, 0xba, 0x0e, 0x18,
	0xe7, 0xd1, 0x82, 0x6d, 0xd2, 0x81, 0x7b, 0x0c, 0xe6, 0x79, 0xd1, 0xd7, 0x0d, 0xa6, 0x28, 0x67,
	0x46, 0xad, 0x67, 0xc6, 0x7d, 0x09, 0xb6, 0x34, 0x17, 0xf8, 0x7e, 0xb9, 0x86, 0x6f, 0xb7, 0x4a,
	0x45, 0xda, 0x48, 0x74, 0x2f, 0x40, 0xf7, 0xf9, 0x60, 0x3d, 0x92, 0x8b, 0x1a, 0xbc, 0xa1, 0xad,
	0xf1, 0x86, 0xfb, 0x1c, 0x2c, 0xf1, 0x90, 0x88, 0xed, 0xae, 0xc5, 0xee, 0xd4, 0x30, 0x70, 0x0b,
	0x19, 0xf9, 0x0d, 0xe8, 0x7e, 0x72, 0x43, 0x73, 0x5e, 0x0b, 0x5b, 0xa5, 0xd5, 0xfc, 0xf3, 0xf3,
	0xb6, 0xfa, 0xd0, 0x67, 0x60, 0xb1, 0x28, 0x26, 0x39, 0xc3, 0x71, 0x2a, 0xa3, 0xd7, 0x17, 0xee,
	0x05, 0xa8, 0x7e, 0xca, 0xf3, 0xc3, 0x61, 0x98, 0x91, 0x3c, 0x97, 0xcf, 0x95, 0x22, 0x3a, 0x02,
	0x23, 0x5f, 0xbe, 0x4f, 0x08, 0x93, 0xd3, 0x7a, 0x50, 0x25, 0x75, 0x29, 0xae, 0x03, 0xa9, 0x76,
	0xff, 0x54, 0xa0, 0x5d, 0x32, 0x0a, 0x3a, 0x05, 0x93, 0x24, 0x2c, 0x8b, 0x48, 0x2e, 0x6b, 0xf9,
	0x7c, 0x83, 0x75, 0x86, 0x5e, 0x61, 0x50, 0xac, 0x50, 0x69, 0xde, 0x7b, 0x05, 0x1f, 0x36, 0x15,
	0x5b, 0xd6, 0xe7, 0xb0, 0xb9, 0x3e, 0x56, 0x73, 0x57, 0x18, 0x98, 0x33, 0xc2, 0xee, 0x69, 0x76,
	0xfb, 0x98, 0xc6, 0xf3, 0xb1, 0xe7, 0x44, 0x9b, 0x25, 0x78, 0x21, 0x70, 0x69, 0x07, 0x95, 0x8c,
	0xbe, 0x06, 0xb3, 0xa8, 0xab, 0x24, 0x83, 0x8d, 0xba, 0x4b, 0x3d, 0x9f, 0x1f, 0x19, 0x75, 0xef,
	0xfc, 0x48, 0x1b, 0xd9, 0xc5, 0x1f, 0xb9, 0xd3, 0x7d, 0x85, 0xd7, 0xc3, 0x74, 0x0b, 0xc2, 0x56,
	0xf7, 0x10, 0xb6, 0x7b, 0x02, 0x46, 0x91, 0x09, 0xaf, 0xed, 0x3a, 0x0a, 0xb3, 0x72, 0x10, 0xf8,
	0x99, 0x37, 0x74, 0x8e, 0x19, 0xb9, 0xc7, 0xab, 0x92, 0x07, 0xa5, 0xe8, 0x3e, 0x05, 0xfd, 0x2d,
	0xbd, 0x25, 0x49, 0x8d, 0xa3, 0xd2, 0xc0, 0x71, 0xf4, 0x77, 0x0b, 0xb4, 0x5f, 0x68, 0x8a, 0xbe,
	0x03, 0xa8, 0x18, 0x25, 0x47, 0x75, 0xf5, 0xc5, 0xce, 0xf5, 0x3e, 0xdd, 0xe4, 0x1d, 0x51, 0xfd,
	0x0b, 0x68, 0xcb, 0x45, 0xd9, 0xe2, 0x74, 0xf8, 0x70, 0x99, 0x84, 0xcb, 0x31, 0x18, 0x62, 0xbe,
	0xb7, 0x38, 0xa0, 0xf5, 0x0d, 0x28, 0x23, 0x48, 0x28, 0xf7, 0x46, 0x68, 0xb6, 0xe4, 0x04, 0x8c,
	0x89, 0x58, 0x37, 0xe4, 0xec, 0xa2, 0xf4, 0xde, 0x16, 0x22, 0xe5, 0xff, 0xe9, 0x67, 0x64, 0x41,
	0x18, 0xd9, 0x56, 0xff, 0x9c, 0xd2, 0xf9, 0x82, 0x14, 0xdf, 0x27, 0xef, 0x97, 0x37, 0x43, 0x8f,
	0x7f, 0xae, 0xa0, 0xd3, 0xc6, 0x12, 0x34, 0x93, 0xa9, 0x5a, 0xbd, 0xd3, 0xf3, 0x1b, 0xd0, 0x2f,
	0x19, 0xce, 0xd8, 0x3e, 0x14, 0xea, 0xd4, 0x9e, 0x43, 0xeb, 0x92, 0xd1, 0xf4, 0xf1, 0x89, 0x1d,
	0x83, 0x3e, 0xa5, 0xf3, 0x28, 0x69, 0x64, 0xd5, 0xf8, 0x7b, 0xe8, 0xd5, 0x5c, 0x53, 0x0c, 0xc7,
	0x29, 0x18, 0x53, 0x3a, 0xa7, 0x4b, 0x86, 0x76, 0x3c, 0xb8, 0x33, 0xd0, 0x10, 0x5a, 0x9c, 0x9f,
	0x76, 0xfa, 0x35, 0x58, 0x8d, 0xd3, 0xd8, 0x6b, 0xf3, 0x77, 0xbd, 0xb0, 0x30, 0xc4, 0xcf, 0xcb,
	0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x59, 0x88, 0x40, 0xbf, 0x0d, 0x0a, 0x00, 0x00,
}
