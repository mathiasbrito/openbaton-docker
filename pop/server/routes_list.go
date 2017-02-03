package server

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

// Container fetches created or running containers from Docker, applying the given filter.
func (svc *service) Containers(ctx context.Context, filter *pop.Filter) (*pop.ContainerList, error) {
	// filter for a container with the given ID
	if filter.Id != "" {
		cont, err := svc.getSingleContainerInfo(ctx, filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.ContainerList{
			List: []*pop.Container{cont},
		}, nil
	}

	return svc.getContainerInfos(ctx)
}

// dummy flavours
var (
	dockerFlavour = &pop.Flavour{
		Id:   "docker-flavour-id",
		Name: "docker.container",
	}
	flavours = &pop.FlavourList{
		List: []*pop.Flavour{dockerFlavour},
	}
)

// Flavours are not necessary; the only reason they are implemented it's because they exist in the
// OpenStack/Amazon/... world, and so the NFVO expects one of them.
// Letting the PoP declare fake containers gives an appearance of continuity with the rest of the NFV world.
func (*service) Flavours(ctx context.Context, filter *pop.Filter) (*pop.FlavourList, error) {
	if filter.Id != "" {
		for _, fl := range flavours.List {
			if fl.Id == filter.Id {
				return &pop.FlavourList{List: []*pop.Flavour{fl}}, nil
			}
		}

		return nil, fmt.Errorf("unsupported flavour with id %s", filter.Id)
	}

	return flavours, nil
}

// Images retrieves and returns the available images on the Docker daemon.
func (svc *service) Images(ctx context.Context, filter *pop.Filter) (*pop.ImageList, error) {
	// filter for an image with the given ID
	if filter.Id != "" {
		img, err := svc.getSingleImageInfo(ctx, filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.ImageList{
			List: []*pop.Image{img},
		}, nil
	}

	return svc.getImageInfos(ctx)
}

// Networks retrieves the current daemon networks.
func (svc *service) Networks(ctx context.Context, filter *pop.Filter) (*pop.NetworkList, error) {
	// filter for a network with the given ID
	if filter.Id != "" {
		netw, err := svc.getSingleNetworkInfo(ctx, filter.Id)
		if err != nil {
			return nil, err
		}

		return &pop.NetworkList{
			List: []*pop.Network{netw},
		}, nil
	}

	return svc.getNetworkInfos(ctx)
}

func (svc *service) getContainerInfos(ctx context.Context) (*pop.ContainerList, error) {
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	// preallocate the list
	conts := make([]*pop.Container, 0, len(svc.conts))

	for _, cj := range svc.conts {
		conts = append(conts, cj.Container)
	}

	return &pop.ContainerList{List: conts}, nil
}

func (svc *service) getDockerContainers(ctx context.Context) ([]types.Container, error) {
	filts, err := filters.FromParam(`{"status": {"created": true, "running": true}}`)
	if err != nil {
		return nil, err
	}

	return svc.cln.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filts,
	})
}

func (svc *service) getImageInfos(ctx context.Context) (*pop.ImageList, error) {
	dockerImgs, err := svc.cln.ImageList(ctx, types.ImageListOptions{})
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

func (svc *service) getNetworkInfos(ctx context.Context) (*pop.NetworkList, error) {
	dockerNets, err := svc.cln.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	nets := make([]*pop.Network, len(dockerNets))
	for i, dnet := range dockerNets {
		nets[i] = extractNetwork(dnet)
	}

	return &pop.NetworkList{List: nets}, nil
}

func (svc *service) getSingleContainerInfo(ctx context.Context, id string) (*pop.Container, error) {
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	// preallocate the list

	cj, found := svc.conts[id]
	if !found {
		return nil, ErrNoSuchContainer
	}

	return cj.Container, nil
}

func (svc *service) getSingleImageInfo(ctx context.Context, id string) (*pop.Image, error) {
	dimg, _, err := svc.cln.ImageInspectWithRaw(ctx, id)
	if err != nil {
		return nil, err
	}

	// why is Docker API such a mess?
	created, err := time.Parse(time.RFC3339Nano, dimg.Created)
	if err != nil {
		return nil, pop.InternalErr
	}

	return &pop.Image{
		Id:      dimg.ID,
		Names:   dimg.RepoTags,
		Created: created.Unix(),
	}, nil
}

func (svc *service) getSingleNetworkInfo(ctx context.Context, id string) (*pop.Network, error) {
	dnet, err := svc.cln.NetworkInspect(ctx, id)
	if err != nil {
		return nil, err
	}

	return extractNetwork(dnet), nil
}

func extractEndpoint(endpointSettings *network.EndpointSettings) *pop.Endpoint {
	var ipv4, ipv6 *pop.Ip

	// IPAMConfig may contain pre-allocated IP addresses for a created, but not yet started, container.
	if endpointSettings.IPAddress != "" {
		fullAddr := fmt.Sprintf("%s/%d", endpointSettings.IPAddress, endpointSettings.IPPrefixLen)
		_, ipnet, err := net.ParseCIDR(fullAddr)
		if err != nil {
			panic("should not happen: " + err.Error())
		}

		ipv4 = &pop.Ip{
			Address: endpointSettings.IPAddress,
			Subnet: &pop.Subnet{
				Cidr:    ipnet.String(),
				Gateway: endpointSettings.Gateway,
			},
		}
	} else {
		if endpointSettings.IPAMConfig != nil {
			ipv4 = &pop.Ip{
				Address: endpointSettings.IPAMConfig.IPv4Address,
			}
		}
	}

	if endpointSettings.GlobalIPv6Address != "" {
		fullAddr := fmt.Sprintf("%s/%d", endpointSettings.GlobalIPv6Address, endpointSettings.GlobalIPv6PrefixLen)
		_, ipnet, err := net.ParseCIDR(fullAddr)
		if err != nil {
			panic("should not happen: " + err.Error())
		}

		ipv6 = &pop.Ip{
			Address: endpointSettings.GlobalIPv6Address,
			Subnet: &pop.Subnet{
				Cidr:    ipnet.String(),
				Gateway: endpointSettings.IPv6Gateway,
			},
		}
	} else {
		if endpointSettings.IPAMConfig != nil {
			ipv6 = &pop.Ip{
				Address: endpointSettings.IPAMConfig.IPv6Address,
			}
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
