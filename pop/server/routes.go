package server

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/openbaton-docker/pop"
)

//go:generate protoc -I ../proto ../proto/pop.proto --go_out=plugins=grpc:..

const (
	// TokenBytes specifies how long a token is.
	TokenBytes = 32

	loginMethod = "/pop.PoP/Login"
)

var (
	AuthErr         = grpc.Errorf(codes.PermissionDenied, "access denied")
	InternalErr     = grpc.Errorf(codes.Internal, "server fault")
	InvalidArgErr   = grpc.Errorf(codes.InvalidArgument, "invalid arguments")
	InvalidTokenErr = grpc.Errorf(codes.PermissionDenied, "invalid token")
	NotLoggedErr    = grpc.Errorf(codes.Unauthenticated, "not authenticated")
)

type service struct {
	sessionManager
	users Users
	cln   *client.Client
}

func newService(cfg Config) (*service, error) {
	cln, err := dialDocker(cfg)
	if err != nil {
		return nil, err
	}

	return &service{
		cln: cln,
		sessionManager: sessionManager{
			tk: make(map[string]struct{}),
		},
		users: cfg.Users,
	}, nil
}

type sessionManager struct {
	l  sync.RWMutex
	tk map[string]struct{}
}

func (sm *sessionManager) CheckToken(tok string) bool {
	sm.l.RLock()
	defer sm.l.RUnlock()

	_, ok := sm.tk[tok]
	return ok
}

func (sm *sessionManager) DeleteToken(tok string) {
	delete(sm.tk, tok)
}

func (sm *sessionManager) NewToken() (string, error) {
	b := make([]byte, TokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)

	sm.l.Lock()
	defer sm.l.Unlock()

	sm.tk[token] = struct{}{}

	return token, nil
}

func (svc *service) Containers(ctx context.Context, filter *pop.Filter) (*pop.ContainerList, error) {
	// filter for a container with the given ID
	if filter.Id != "" {
		cont, err := svc.getSingleContainerInfo(filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.ContainerList{
			List: []*pop.Container{cont},
		}, nil
	}

	return svc.getContainerInfos()
}

func (svc *service) Images(ctx context.Context, filter *pop.Filter) (*pop.ImageList, error) {
	// filter for an image with the given ID
	if filter.Id != "" {
		img, err := svc.getSingleImageInfo(filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.ImageList{
			List: []*pop.Image{img},
		}, nil
	}

	return svc.getImageInfos()
}

// Login logs into the Pop. It should always be the first function called (to setup a token).
// Remember that tokens are transient and not stored, so a new login is needed in case the service dies.
func (svc *service) Login(ctx context.Context, creds *pop.Credentials) (*pop.Token, error) {
	if creds == nil {
		return nil, InvalidArgErr
	}

	if user, found := svc.users[creds.Username]; found {
		if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(creds.Password)) == nil {
			tok, err := svc.NewToken()
			if err != nil {
				return nil, InternalErr
			}

			return &pop.Token{Value: tok}, nil
		}
	}

	return nil, AuthErr
}

func (svc *service) Logout(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	// getToken() will always return a valid token (it has been checked in unaryInterceptor()).

	svc.DeleteToken(getToken(ctx))

	return &empty.Empty{}, nil
}

func (svc *service) Networks(ctx context.Context, filter *pop.Filter) (*pop.NetworkList, error) {
	// filter for a network with the given ID
	if filter.Id != "" {
		netw, err := svc.getSingleNetworkInfo(filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.NetworkList{
			List: []*pop.Network{netw},
		}, nil
	}

	return svc.getNetworkInfos()
}

func (svc *service) authorize(ctx context.Context) error {
	token := getToken(ctx)
	if token == "" {
		return NotLoggedErr
	}

	if !svc.CheckToken(token) {
		return InvalidTokenErr
	}

	return nil
}

func (svc *service) getContainerInfos() (*pop.ContainerList, error) {
	dockerConts, err := svc.getDockerContainersForStatus("created")
	if err != nil {
		return nil, err
	}

	runningConts, err := svc.getDockerContainersForStatus("running")
	if err != nil {
		return nil, err
	}

	dockerConts = append(dockerConts, runningConts...)

	conts := make([]*pop.Container, len(dockerConts))

	for i, dcont := range dockerConts {
		conts[i] = &pop.Container{
			Id:             dcont.ID,
			Names:          dcont.Names,
			Status:         dcont.State,
			ExtendedStatus: dcont.Status, // The Docker API is not very clear about this
			ImageId:        dcont.ImageID,
			Created:        dcont.Created,
			Command:        dcont.Command,
			Endpoints:      extractEndpoints(dcont.NetworkSettings.Networks),
		}
	}

	return &pop.ContainerList{List: conts}, nil
}

func (svc *service) getDockerContainersForStatus(status string) ([]types.Container, error) {
	filts, err := filters.FromParam("status=" + status)
	if err != nil {
		return nil, err
	}

	return svc.cln.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: filts,
	})
}

func (svc *service) getImageInfos() (*pop.ImageList, error) {
	dockerImgs, err := svc.cln.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	imgs := make([]*pop.Image, len(dockerImgs))
	for i, dimg := range dockerImgs {
		imgs[i] = &pop.Image{
			Id:      dimg.ID,
			Names:   dimg.RepoTags,
			Created: dimg.Created,
		}
	}

	return &pop.ImageList{List: imgs}, nil
}

func (svc *service) getNetworkInfos() (*pop.NetworkList, error) {
	dockerNets, err := svc.cln.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	nets := make([]*pop.Network, len(dockerNets))
	for i, dnet := range dockerNets {
		nets[i] = extractNetwork(dnet)
	}

	return &pop.NetworkList{List: nets}, nil
}

func (svc *service) getSingleContainerInfo(id string) (*pop.Container, error) {
	dcont, err := svc.cln.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	// why is Docker API such a mess?
	created, err := time.Parse(time.RFC3339Nano, dcont.Created)
	if err != nil {
		return nil, InternalErr
	}

	b := bytes.Buffer{}
	for _, part := range dcont.Config.Cmd {
		b.WriteString(part)
		b.WriteRune(' ')
	}

	return &pop.Container{
		Id:             dcont.ID,
		Names:          []string{dcont.Name},
		Status:         dcont.State.Status,
		ExtendedStatus: dcont.State.Error,
		ImageId:        dcont.Image,
		Created:        created.Unix(),
		Command:        b.String(),
		Endpoints:      extractEndpoints(dcont.NetworkSettings.Networks),
	}, nil
}

func (svc *service) getSingleImageInfo(id string) (*pop.Image, error) {
	dimg, _, err := svc.cln.ImageInspectWithRaw(context.Background(), id)
	if err != nil {
		return nil, err
	}

	// why is Docker API such a mess?
	created, err := time.Parse(time.RFC3339Nano, dimg.Created)
	if err != nil {
		return nil, InternalErr
	}

	return &pop.Image{
		Id:      dimg.ID,
		Names:   dimg.RepoTags,
		Created: created.Unix(),
	}, nil
}

func (svc *service) getSingleNetworkInfo(id string) (*pop.Network, error) {
	dnet, err := svc.cln.NetworkInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return extractNetwork(dnet), nil
}

func (svc *service) streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := svc.authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

func (svc *service) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Let the Login method AND ONLY IT pass through without a valid token (for obvious reasons)
	if info.FullMethod != loginMethod {
		if err := svc.authorize(ctx); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func extractEndpoint(endpointSettings *network.EndpointSettings) *pop.Endpoint {
	var ipv4, ipv6 *pop.Ip

	// IPAMConfig may contain pre-allocated IP addresses for a created, but not yet started, container.
	if endpointSettings.IPAddress != "" {
		if endpointSettings.IPAMConfig != nil {
			ipv4 = &pop.Ip{
				Address: endpointSettings.IPAMConfig.IPv4Address,
			}
		}
	} else {
		fullAddr := fmt.Sprintf("%s/%d", endpointSettings.IPAddress, endpointSettings.IPPrefixLen)
		_, ipnet, err := net.ParseCIDR(fullAddr)
		if err != nil {
			panic("should not happen")
		}

		ipv4 = &pop.Ip{
			Address: endpointSettings.IPAddress,
			Subnet: &pop.Subnet{
				Cidr:    ipnet.String(),
				Gateway: endpointSettings.Gateway,
			},
		}
	}

	if endpointSettings.GlobalIPv6Address != "" {
		if endpointSettings.IPAMConfig != nil {
			ipv6 = &pop.Ip{
				Address: endpointSettings.IPAMConfig.IPv6Address,
			}
		}
	} else {
		fullAddr := fmt.Sprintf("%s/%d", endpointSettings.GlobalIPv6Address, endpointSettings.GlobalIPv6PrefixLen)
		_, ipnet, err := net.ParseCIDR(fullAddr)
		if err != nil {
			panic("should not happen")
		}

		ipv6 = &pop.Ip{
			Address: endpointSettings.GlobalIPv6Address,
			Subnet: &pop.Subnet{
				Cidr:    ipnet.String(),
				Gateway: endpointSettings.IPv6Gateway,
			},
		}
	}

	return &pop.Endpoint{
		NetId:      endpointSettings.NetworkID,
		EndpointId: endpointSettings.EndpointID,
		Ipv4:       ipv4,
		Ipv6:       ipv6,
	}
}

func extractEndpoints(dNetMap map[string]*network.EndpointSettings) map[string]*pop.Endpoint {
	endpoints := make(map[string]*pop.Endpoint)

	for netname, endpointSettings := range dNetMap {
		endpoints[netname] = extractEndpoint(endpointSettings)
	}

	return endpoints
}

func extractNetwork(dnet types.NetworkResource) *pop.Network {
	subs := extractSubnets(dnet.IPAM.Config)

	return &pop.Network{
		Id:       dnet.ID,
		Name:     dnet.Name,
		External: !dnet.Internal,
		Subnets:  subs,
	}
}

func extractSubnets(dSubnets []network.IPAMConfig) []*pop.Subnet {
	subs := make([]*pop.Subnet, len(dSubnets))

	for i, dSubnet := range dSubnets {
		subs[i] = &pop.Subnet{
			Cidr:    dSubnet.Subnet,
			Gateway: dSubnet.Gateway,
		}
	}

	return subs
}

func getToken(ctx context.Context) string {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return ""
	}

	if len(md["token"]) == 0 {
		return ""
	}

	return md["token"][0]
}

func dialDocker(cfg Config) (*client.Client, error) {
	return client.NewClient(cfg.DockerdHost, client.DefaultVersion, nil, nil)
}
