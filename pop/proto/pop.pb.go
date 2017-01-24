// Code generated by protoc-gen-go.
// source: pop.proto
// DO NOT EDIT!

/*
Package pop is a generated protocol buffer package.

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
	Network
	NetworkList
	Subnet
	Token
*/
package pop

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

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

type Container struct {
	Id             string               `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Names          []string             `protobuf:"bytes,2,rep,name=names" json:"names,omitempty"`
	ImageId        string               `protobuf:"bytes,3,opt,name=image_id,json=imageId" json:"image_id,omitempty"`
	Command        string               `protobuf:"bytes,4,opt,name=command" json:"command,omitempty"`
	Created        int64                `protobuf:"varint,5,opt,name=created" json:"created,omitempty"`
	Status         string               `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	ExtendedStatus string               `protobuf:"bytes,7,opt,name=extended_status,json=extendedStatus" json:"extended_status,omitempty"`
	Endpoints      map[string]*Endpoint `protobuf:"bytes,8,rep,name=endpoints" json:"endpoints,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Container) Reset()                    { *m = Container{} }
func (m *Container) String() string            { return proto.CompactTextString(m) }
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

func (m *Container) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
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

type ContainerConfig struct {
	Name      string               `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ImageId   string               `protobuf:"bytes,2,opt,name=image_id,json=imageId" json:"image_id,omitempty"`
	FlavourId string               `protobuf:"bytes,3,opt,name=flavour_id,json=flavourId" json:"flavour_id,omitempty"`
	Endpoints map[string]*Endpoint `protobuf:"bytes,4,rep,name=endpoints" json:"endpoints,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ContainerConfig) Reset()                    { *m = ContainerConfig{} }
func (m *ContainerConfig) String() string            { return proto.CompactTextString(m) }
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
func (m *ContainerList) String() string            { return proto.CompactTextString(m) }
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
func (m *Credentials) String() string            { return proto.CompactTextString(m) }
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
func (m *Endpoint) String() string            { return proto.CompactTextString(m) }
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
func (m *Filter) String() string            { return proto.CompactTextString(m) }
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
func (m *Flavour) String() string            { return proto.CompactTextString(m) }
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
func (m *FlavourList) String() string            { return proto.CompactTextString(m) }
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
func (m *Image) String() string            { return proto.CompactTextString(m) }
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
func (m *ImageList) String() string            { return proto.CompactTextString(m) }
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
func (m *Infos) String() string            { return proto.CompactTextString(m) }
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
func (m *Ip) String() string            { return proto.CompactTextString(m) }
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

type Network struct {
	Id       string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name     string    `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	External bool      `protobuf:"varint,3,opt,name=external" json:"external,omitempty"`
	Subnets  []*Subnet `protobuf:"bytes,4,rep,name=subnets" json:"subnets,omitempty"`
}

func (m *Network) Reset()                    { *m = Network{} }
func (m *Network) String() string            { return proto.CompactTextString(m) }
func (*Network) ProtoMessage()               {}
func (*Network) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

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
func (m *NetworkList) String() string            { return proto.CompactTextString(m) }
func (*NetworkList) ProtoMessage()               {}
func (*NetworkList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *NetworkList) GetList() []*Network {
	if m != nil {
		return m.List
	}
	return nil
}

type Subnet struct {
	Cidr    string `protobuf:"bytes,1,opt,name=cidr" json:"cidr,omitempty"`
	Gateway string `protobuf:"bytes,2,opt,name=gateway" json:"gateway,omitempty"`
}

func (m *Subnet) Reset()                    { *m = Subnet{} }
func (m *Subnet) String() string            { return proto.CompactTextString(m) }
func (*Subnet) ProtoMessage()               {}
func (*Subnet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

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
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *Token) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*Container)(nil), "pop.Container")
	proto.RegisterType((*ContainerConfig)(nil), "pop.ContainerConfig")
	proto.RegisterType((*ContainerList)(nil), "pop.ContainerList")
	proto.RegisterType((*Credentials)(nil), "pop.Credentials")
	proto.RegisterType((*Endpoint)(nil), "pop.Endpoint")
	proto.RegisterType((*Filter)(nil), "pop.Filter")
	proto.RegisterType((*Flavour)(nil), "pop.Flavour")
	proto.RegisterType((*FlavourList)(nil), "pop.FlavourList")
	proto.RegisterType((*Image)(nil), "pop.Image")
	proto.RegisterType((*ImageList)(nil), "pop.ImageList")
	proto.RegisterType((*Infos)(nil), "pop.Infos")
	proto.RegisterType((*Ip)(nil), "pop.Ip")
	proto.RegisterType((*Network)(nil), "pop.Network")
	proto.RegisterType((*NetworkList)(nil), "pop.NetworkList")
	proto.RegisterType((*Subnet)(nil), "pop.Subnet")
	proto.RegisterType((*Token)(nil), "pop.Token")
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
	// Delete deletes the container identified by the given filter.
	Delete(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	// Start starts the container identified by the given filter.
	Start(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Container, error)
	// Stop stops the container identified by the given filter.
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
	err := grpc.Invoke(ctx, "/pop.Pop/Containers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Flavours(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*FlavourList, error) {
	out := new(FlavourList)
	err := grpc.Invoke(ctx, "/pop.Pop/Flavours", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Images(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*ImageList, error) {
	out := new(ImageList)
	err := grpc.Invoke(ctx, "/pop.Pop/Images", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Networks(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*NetworkList, error) {
	out := new(NetworkList)
	err := grpc.Invoke(ctx, "/pop.Pop/Networks", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Create(ctx context.Context, in *ContainerConfig, opts ...grpc.CallOption) (*Container, error) {
	out := new(Container)
	err := grpc.Invoke(ctx, "/pop.Pop/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Delete(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/pop.Pop/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Start(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Container, error) {
	out := new(Container)
	err := grpc.Invoke(ctx, "/pop.Pop/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Stop(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/pop.Pop/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Login(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := grpc.Invoke(ctx, "/pop.Pop/Login", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Logout(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/pop.Pop/Logout", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *popClient) Info(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*Infos, error) {
	out := new(Infos)
	err := grpc.Invoke(ctx, "/pop.Pop/Info", in, out, c.cc, opts...)
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
	// Delete deletes the container identified by the given filter.
	Delete(context.Context, *Filter) (*google_protobuf.Empty, error)
	// Start starts the container identified by the given filter.
	Start(context.Context, *Filter) (*Container, error)
	// Stop stops the container identified by the given filter.
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
		FullMethod: "/pop.Pop/Containers",
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
		FullMethod: "/pop.Pop/Flavours",
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
		FullMethod: "/pop.Pop/Images",
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
		FullMethod: "/pop.Pop/Networks",
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
		FullMethod: "/pop.Pop/Create",
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
		FullMethod: "/pop.Pop/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Delete(ctx, req.(*Filter))
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
		FullMethod: "/pop.Pop/Start",
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
		FullMethod: "/pop.Pop/Stop",
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
		FullMethod: "/pop.Pop/Login",
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
		FullMethod: "/pop.Pop/Logout",
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
		FullMethod: "/pop.Pop/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PopServer).Info(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Pop_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pop.Pop",
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

func init() { proto.RegisterFile("pop.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 817 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x55, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x96, 0xed, 0xd8, 0x89, 0x8f, 0xd9, 0xec, 0x6a, 0xb4, 0xac, 0x4c, 0xa0, 0x10, 0xb9, 0x54,
	0x04, 0x56, 0x0d, 0x52, 0x17, 0x55, 0x2b, 0xb8, 0x42, 0xa1, 0x8b, 0x2c, 0x0a, 0x42, 0x2e, 0xf7,
	0xd5, 0xb4, 0x33, 0x8d, 0xac, 0xda, 0x33, 0x23, 0xcf, 0xa4, 0x25, 0x2f, 0xc1, 0x8b, 0xf1, 0x30,
	0xf0, 0x08, 0x68, 0x7e, 0xec, 0xc4, 0x69, 0x22, 0xf5, 0x8a, 0xbb, 0x39, 0xe7, 0x3b, 0x3e, 0x3f,
	0xdf, 0x7c, 0x67, 0x0c, 0xb1, 0xe0, 0x62, 0x2e, 0x1a, 0xae, 0x38, 0x0a, 0x04, 0x17, 0x93, 0x84,
	0xd6, 0x42, 0xad, 0xad, 0x27, 0xfb, 0xdb, 0x87, 0x78, 0xc1, 0x99, 0xc2, 0x25, 0xa3, 0x0d, 0x1a,
	0x83, 0x5f, 0x92, 0xd4, 0x9b, 0x7a, 0xb3, 0xb8, 0xf0, 0x4b, 0x82, 0x5e, 0x43, 0xc8, 0x70, 0x4d,
	0x65, 0xea, 0x4f, 0x83, 0x59, 0x5c, 0x58, 0x03, 0x7d, 0x02, 0xa3, 0xb2, 0xc6, 0x4b, 0x7a, 0x5d,
	0x92, 0x34, 0x30, 0xb1, 0x43, 0x63, 0xe7, 0x04, 0xa5, 0x30, 0xbc, 0xe5, 0x75, 0x8d, 0x19, 0x49,
	0x07, 0x16, 0x71, 0xa6, 0x41, 0x1a, 0x8a, 0x15, 0x25, 0x69, 0x38, 0xf5, 0x66, 0x41, 0xd1, 0x9a,
	0xe8, 0x0d, 0x44, 0x52, 0x61, 0xb5, 0x92, 0x69, 0x64, 0x3e, 0x71, 0x16, 0xfa, 0x0a, 0x5e, 0xd2,
	0x3f, 0x15, 0x65, 0x84, 0x92, 0x6b, 0x17, 0x30, 0x34, 0x01, 0xe3, 0xd6, 0x7d, 0x65, 0x03, 0x7f,
	0x80, 0x98, 0x32, 0x22, 0x78, 0xc9, 0x94, 0x4c, 0x47, 0xd3, 0x60, 0x96, 0x9c, 0x1d, 0xcd, 0xf5,
	0xd0, 0xdd, 0x60, 0xf3, 0x8b, 0x16, 0xbf, 0x60, 0xaa, 0x59, 0x17, 0x9b, 0xf8, 0xc9, 0x2f, 0x30,
	0xee, 0x83, 0xe8, 0x15, 0x04, 0xf7, 0x74, 0xed, 0x58, 0xd0, 0x47, 0x74, 0x0c, 0xe1, 0x03, 0xae,
	0x56, 0x34, 0xf5, 0xa7, 0xde, 0x2c, 0x39, 0x7b, 0x61, 0x92, 0xb7, 0x5f, 0x15, 0x16, 0xfb, 0xde,
	0x7f, 0xef, 0x65, 0xff, 0x78, 0xf0, 0xb2, 0x2b, 0xba, 0xe0, 0xec, 0xae, 0x5c, 0x22, 0x04, 0x03,
	0x4d, 0x9b, 0xcb, 0x67, 0xce, 0x3d, 0x06, 0xfd, 0x3e, 0x83, 0x47, 0x00, 0x77, 0x15, 0x7e, 0xe0,
	0xab, 0x66, 0x43, 0x6f, 0xec, 0x3c, 0x39, 0x41, 0x3f, 0x6e, 0xcf, 0x3a, 0x30, 0xb3, 0x1e, 0xf7,
	0x67, 0xb5, 0x65, 0xff, 0xaf, 0x89, 0xdf, 0xc1, 0x8b, 0xae, 0xf2, 0x65, 0x29, 0x15, 0xca, 0x60,
	0x50, 0x95, 0x52, 0xa5, 0x9e, 0xe9, 0x6d, 0xdc, 0xef, 0xad, 0x30, 0x58, 0x76, 0x01, 0xc9, 0xa2,
	0xa1, 0x84, 0x32, 0x55, 0xe2, 0x4a, 0xa2, 0x09, 0x8c, 0x56, 0x92, 0x36, 0x5b, 0x2c, 0x75, 0xb6,
	0xc6, 0x04, 0x96, 0xf2, 0x91, 0x37, 0x2d, 0x53, 0x9d, 0x9d, 0xfd, 0xe5, 0xc1, 0xa8, 0xed, 0x09,
	0x7d, 0x0c, 0x11, 0xa3, 0xea, 0xba, 0x93, 0x6f, 0xc8, 0xa8, 0xca, 0x09, 0xfa, 0x02, 0x92, 0x76,
	0xf2, 0x0d, 0xd9, 0xd0, 0xba, 0x72, 0x82, 0x3e, 0x85, 0x41, 0x29, 0x1e, 0xbe, 0x33, 0x4c, 0x27,
	0x67, 0x43, 0xd3, 0x6f, 0x2e, 0x0a, 0xe3, 0x74, 0xe0, 0xb9, 0xd1, 0xf2, 0x0e, 0x78, 0xae, 0x59,
	0xab, 0xf1, 0xad, 0x51, 0x73, 0x5c, 0xe8, 0x63, 0x96, 0x42, 0xf4, 0xa1, 0xac, 0xd4, 0xd3, 0x45,
	0xca, 0x4e, 0x61, 0xf8, 0xc1, 0xde, 0xe1, 0x93, 0x1d, 0x6b, 0xf5, 0xe1, 0x6f, 0xf4, 0x91, 0x7d,
	0x0b, 0x89, 0x0b, 0x37, 0x9c, 0x4e, 0x7b, 0x9c, 0x7e, 0x64, 0xda, 0x70, 0xb8, 0x63, 0xf4, 0x67,
	0x08, 0x73, 0x2d, 0xa0, 0x67, 0x6e, 0xf0, 0xd6, 0x32, 0x06, 0xbd, 0x65, 0xcc, 0xde, 0x42, 0x6c,
	0x12, 0x99, 0xba, 0x9f, 0xf7, 0xea, 0x82, 0x1d, 0x5f, 0xa3, 0xae, 0xea, 0xaf, 0x10, 0xe6, 0xec,
	0x8e, 0x4b, 0x3d, 0x83, 0x5a, 0x8b, 0x4e, 0xe3, 0xfa, 0xbc, 0x6f, 0x2e, 0xf4, 0x19, 0xc4, 0xaa,
	0xac, 0xa9, 0x54, 0xb8, 0x16, 0xae, 0xf2, 0xc6, 0x91, 0x2d, 0xc0, 0xcf, 0x85, 0xee, 0x0d, 0x13,
	0xd2, 0x50, 0x29, 0x5d, 0xba, 0xd6, 0x44, 0xc7, 0x10, 0xc9, 0xd5, 0x0d, 0xa3, 0xca, 0xa9, 0x32,
	0x31, 0x0d, 0x5d, 0x19, 0x57, 0xe1, 0xa0, 0x4c, 0xc0, 0xf0, 0x37, 0xaa, 0x1e, 0x79, 0x73, 0xff,
	0x1c, 0xa6, 0xb5, 0xbe, 0xf4, 0x6b, 0xd2, 0x30, 0x5c, 0x99, 0x86, 0x46, 0x45, 0x67, 0xa3, 0x13,
	0x18, 0xda, 0xa4, 0xed, 0xa6, 0xf5, 0x0a, 0xb6, 0x98, 0xbe, 0x2c, 0x57, 0xf1, 0xe0, 0x65, 0x39,
	0xdc, 0xd1, 0x76, 0x0e, 0x91, 0xcd, 0xa1, 0x3b, 0xba, 0x2d, 0x49, 0xd3, 0xf2, 0xa6, 0xcf, 0x7a,
	0xfe, 0x25, 0x56, 0xf4, 0x11, 0xaf, 0xdb, 0xa7, 0xc1, 0x99, 0xd9, 0x11, 0x84, 0x7f, 0xf0, 0x7b,
	0xca, 0xf4, 0xa5, 0xda, 0xed, 0x74, 0x52, 0x37, 0xc6, 0xd9, 0xbf, 0x01, 0x04, 0xbf, 0x73, 0x81,
	0x4e, 0x01, 0xba, 0x85, 0x93, 0xc8, 0xf6, 0x6c, 0x65, 0x39, 0x41, 0xfd, 0x75, 0x34, 0xfd, 0x7e,
	0x0d, 0x23, 0xa7, 0xa5, 0x9d, 0xe0, 0x57, 0xdb, 0x3a, 0x33, 0xa1, 0x27, 0x10, 0x99, 0xeb, 0xdf,
	0x09, 0x1c, 0x6f, 0x84, 0xd1, 0x66, 0x74, 0x03, 0xef, 0xcd, 0xb8, 0x4d, 0xd6, 0x1c, 0xa2, 0x85,
	0x51, 0x1e, 0x7a, 0xbd, 0xef, 0x15, 0x9b, 0xec, 0xbc, 0x1f, 0xe8, 0x14, 0xa2, 0x9f, 0x68, 0x45,
	0x15, 0xed, 0x27, 0x7e, 0x33, 0x5f, 0x72, 0xbe, 0xac, 0xa8, 0xfd, 0xa9, 0xdd, 0xac, 0xee, 0xe6,
	0x17, 0xfa, 0x1f, 0x87, 0xbe, 0x84, 0xf0, 0x4a, 0xe1, 0x46, 0xed, 0xeb, 0x77, 0x93, 0xf4, 0x2d,
	0x0c, 0xae, 0x14, 0x17, 0xcf, 0x4b, 0x79, 0x02, 0xe1, 0x25, 0x5f, 0x96, 0x0c, 0xd9, 0x61, 0xb6,
	0xde, 0xb1, 0x89, 0x5d, 0x10, 0x7b, 0x45, 0xef, 0x21, 0xba, 0xe4, 0x4b, 0xbe, 0x52, 0xe8, 0x40,
	0xa2, 0x83, 0x05, 0xbe, 0x81, 0x81, 0x5e, 0xaa, 0x83, 0xdf, 0xb9, 0x35, 0xd4, 0x7b, 0x77, 0x13,
	0x19, 0xec, 0xdd, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x75, 0x7e, 0xb2, 0xf1, 0xe3, 0x07, 0x00,
	0x00,
}
