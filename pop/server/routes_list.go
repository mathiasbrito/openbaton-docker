package server

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

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
		return nil, pop.InternalErr
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
		return nil, pop.InternalErr
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
