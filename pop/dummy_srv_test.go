package pop_test

// Dummy testing server, unsafe and badly written for obvious reasons

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	grpc_md "google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	"fmt"
	"google.golang.org/grpc/codes"
)

var (
	tokens = make(map[string]struct{})

	nets = []*pop.Network{
		{
			Id: util.GenerateID(),
			Name: "private",
			External: false,
			Subnets: []*pop.Subnet{
				{
					Cidr: "172.18.0.0/16",
					Gateway: "172.18.0.1",
				},
			},
		},
	}

	flavs = []*pop.Flavour{
		{
			Id:   "dummy-flavour-id",
			Name: "dummy.container",
		},
	}

	imgs = []*pop.Image{
		{
			Id: util.GenerateID(),
			Names: []string{"nginx:latest"},
			Created: time.Now().Unix(),
		},
	}

	conts = []*pop.Container{
		{
			Id: util.GenerateID(),
			Names: []string{"cont1"},
			ImageId: imgs[0].Id,
			FlavourId: flavs[0].Id,
			Created: time.Now().Unix() - 30,
			Started: time.Now().Unix() - 29, 
			Status: pop.Container_RUNNING,
			ExtendedStatus: "dummy is obviously ok",
			Endpoints: map[string]*pop.Endpoint{
				"private": &pop.Endpoint{
					NetId: nets[0].Id,
					NetName: nets[0].Name,
					Mac: "0e:f6:ae:ab:a5:5a",
					Ipv4: &pop.Ip{
						Address: "172.18.0.22",
						Subnet: nets[0].Subnets[0],
					},
				},
			},
			Md: &pop.Metadata{Entries: make(map[string]string)},
		},
	}
)

// Dummy service
type server struct {
	name  string
	quitChan chan struct{}
}

func dummyServer() (*server, error) {
	return &server{
		name: "dummy",
		quitChan:  make(chan struct{}),
	}, nil
}

func (srv *server) Serve() error {
	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}

	gsrv := grpc.NewServer(
		grpc.StreamInterceptor(srv.streamInterceptor),
		grpc.UnaryInterceptor(srv.unaryInterceptor),
	)

	pop.RegisterPopServer(gsrv, srv)

	return gsrv.Serve(lis)
}

func (srv *server) Containers(ctx context.Context, filter *pop.Filter) (*pop.ContainerList, error) {
	// filter for a container with the given ID
	if filter.Options != nil {
		cont, err := srv.getSingleContainerInfo(ctx, filter)
		if err != nil {
			return nil, err
		}

		return &pop.ContainerList{
			List: []*pop.Container{cont},
		}, nil
	}

	return &pop.ContainerList{
		List: conts,
	}, nil
}

// Flavours are not necessary; the only reason they are implemented it's because they exist in the
// OpenStack/Amazon/... world, and so the NFVO expects one of them.
// We're letting the PoP declare fake flavours, giving an appearance of continuity with the rest of the NFV world.
func (srv *server) Flavours(ctx context.Context, filter *pop.Filter) (*pop.FlavourList, error) {
	if filter.Options != nil {
		fl, err := srv.getSingleFlavourInfo(ctx, filter)
		if err != nil {

			return nil, err
		}

		return &pop.FlavourList{
			List: []*pop.Flavour{fl},
		}, nil
	}

	return &pop.FlavourList{
		List: flavs,
	}, nil
}

// Images retrieves and returns the available images on the Docker daemon.
func (srv *server) Images(ctx context.Context, filter *pop.Filter) (*pop.ImageList, error) {
	// filter for an image with the given ID or name
	if filter.Options != nil {
		img, err := srv.getSingleImageInfo(ctx, filter)
		if err != nil {
			return nil, err
		}

		return &pop.ImageList{
			List: []*pop.Image{img},
		}, nil
	}

	return &pop.ImageList{
		List: imgs,
	}, nil
}

func (srv *server) Login(ctx context.Context, creds *pop.Credentials) (*pop.Token, error) {
	if creds == nil {
		return nil, pop.InvalidArgErr
	}

	if creds.Username == uname && creds.Password == pass {
		tk := util.GenerateID()

		tokens[tk] = struct{}{}

		return &pop.Token{Value: tk}, nil
	}

	return nil, pop.AuthErr
}

func (srv *server) Logout(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	for _, token := range getTokens(ctx) {
		delete(tokens, token)
	}

	return &empty.Empty{}, nil
}

// Networks retrieves the current daemon networks.
func (srv *server) Networks(ctx context.Context, filter *pop.Filter) (*pop.NetworkList, error) {

	// filter for a network with the given ID or name.
	if filter.Options != nil {
		netw, err := srv.getSingleNetworkInfo(ctx, filter)
		if err != nil {
			return nil, err
		}

		return &pop.NetworkList{
			List: []*pop.Network{netw},
		}, nil
	}

	return &pop.NetworkList{
		List: nets,
	}, nil
}

func (srv *server) Info(context.Context, *empty.Empty) (*pop.Infos, error) {
	return &pop.Infos{
		Name:      srv.name,
		Type:      "dummy",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (srv *server) Create(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.FlavourId != "" && cfg.FlavourId != flavs[0].Id {
		return nil, fmt.Errorf("unsupported flavour %v, only %v is supported", cfg.FlavourId, flavs[0].Id)
	}

	if _, err := srv.getSingleContainerInfo(ctx, &pop.Filter{Options: &pop.Filter_Name{Name: cfg.Name}}); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "container name %s already taken", cfg.Name)
	}

	cont := &pop.Container{
		Id: util.GenerateID(),
		Created: time.Now().Unix(),
		Status: pop.Container_CREATED,
		Names: []string{cfg.Name},
		Endpoints: cfg.Endpoints,
		Md: &pop.Metadata{Entries: make(map[string]string)},
	}

	conts = append(conts, cont)

	return cont, nil
}

func contains(h []string, n string) bool {
	for _, e := range h {
		if e == n {
			return true
		}
	}

	return false
}

func (srv *server) Delete(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	id := filter.GetId()
	name := filter.GetName()

	for i, cont := range conts {
		if cont.Id == id || contains(cont.Names, name) {
			cont.Status = pop.Container_EXITED
			conts = append(conts[:i], conts[i+1:]...)

			return &empty.Empty{}, nil
		}
	}

	return nil, grpc.Errorf(codes.NotFound, "container not found")
}

func (srv *server) Metadata(ctx context.Context, newMD *pop.NewMetadata) (*empty.Empty, error) {
	id := newMD.Filter.GetId()
	name := newMD.Filter.GetName()

	for _, cont := range conts {
		if cont.Id == id || contains(cont.Names, name) {
			for k, v := range newMD.Md.Entries {
				if v == "" {
					delete(cont.Md.Entries, k)
				} else {
					cont.Md.Entries[k] = v
				}
			}

			return &empty.Empty{}, nil
		}
	}

	return nil, grpc.Errorf(codes.NotFound, "container not found")
}

func (srv *server) Start(ctx context.Context, filter *pop.Filter) (*pop.Container, error) {
	id := filter.GetId()
	name := filter.GetName()

	for _, cont := range conts {
		if cont.Id == id || contains(cont.Names, name) {
			if cont.Status != pop.Container_CREATED {
				return nil, grpc.Errorf(codes.PermissionDenied, "container in wrong state")
			}

			cont.Status = pop.Container_RUNNING
			cont.ExtendedStatus = "the container is running"

			return cont, nil
		}
	}

	return nil, grpc.Errorf(codes.NotFound, "container not found")
}

func (srv *server) Stop(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	id := filter.GetId()
	name := filter.GetName()

	for _, cont := range conts {
		if cont.Id == id || contains(cont.Names, name) {
			if cont.Status != pop.Container_RUNNING {
				return nil, grpc.Errorf(codes.PermissionDenied, "container in wrong state")
			}

			cont.Status = pop.Container_EXITED
			cont.ExtendedStatus = "the container has exited"

			return &empty.Empty{}, nil
		}
	}

	return nil, grpc.Errorf(codes.NotFound, "container not found")
}

func (srv *server) close() error {
	srv.quitChan <- struct{}{}

	select {
	case <-srv.quitChan:
		return nil

	case <-time.After(5 * time.Second):
		return errors.New("timed out while closing the Docker monitor routine")
	}
}

// authorize checks if the current context is autheticated (ie, if it contains a valid token).
func (srv *server) authorize(ctx context.Context) error {
	toks := getTokens(ctx)

	if len(toks) == 0 {
		return pop.NotLoggedErr
	}

	for _, token := range toks {
		if _, found := tokens[token]; found {
			return nil
		}
	}

	return pop.InvalidTokenErr
}

// streamInterceptor is an interceptor for stream requests.
func (srv *server) streamInterceptor(s interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := srv.authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

// unaryInterceptor intercepts every unary request, and ensures that the caller is authorized before doing anything.
func (srv *server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Let the Login method AND ONLY IT pass through without a valid token (for obvious reasons)
	if info.FullMethod != pop.LoginMethod {
		if err := srv.authorize(ctx); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func (srv *server) getSingleContainerInfo(_ context.Context, filter *pop.Filter) (*pop.Container, error) {
	id := filter.GetId()
	name := filter.GetName()

	if id == "" && name == "" {
		return nil, errors.New("no container specified")
	}

	for _, cont := range conts {
		if cont.Id == id {
			return cont, nil
		}

		for _, n := range cont.Names {
			if n == name {
				return cont, nil
			}
		}
	}

	return nil, errors.New("no container found")
}

func (srv *server) getSingleFlavourInfo(_ context.Context, filter *pop.Filter) (*pop.Flavour, error) {
	id := filter.GetId()
	name := filter.GetName()

	if id == "" && name == "" {
		return nil, errors.New("no flavour specified")
	}

	for _, flav := range flavs {
		if flav.Id == id || flav.Name == name {
			return flav, nil
		}
	}

	return nil, errors.New("no flavour found")
}

func (srv *server) getSingleImageInfo(_ context.Context, filter *pop.Filter) (*pop.Image, error) {
	id := filter.GetId()
	name := filter.GetName()

	if id == "" && name == "" {
		return nil, errors.New("no image specified")
	}

	for _, img := range imgs {
		if img.Id == id {
			return img, nil
		}

		for _, n := range img.Names {
			if n == name {
				return img, nil
			}
		}
	}

	return nil, errors.New("no image found")
}

func (srv *server) getSingleNetworkInfo(_ context.Context, filter *pop.Filter) (*pop.Network, error) {
	id := filter.GetId()
	name := filter.GetName()

	if id == "" && name == "" {
		return nil, errors.New("no network specified")
	}

	for _, net := range nets {
		if net.Id == id || net.Name == name {
			return net, nil
		}
	}

	return nil, errors.New("no network found")
}

// getToken retrieves the tokens from the current context metadata.
func getTokens(ctx context.Context) []string {
	md, ok := grpc_md.FromContext(ctx)
	if !ok {
		return []string{}
	}

	if len(md["token"]) == 0 {
		return []string{}
	}

	return md["token"]
}