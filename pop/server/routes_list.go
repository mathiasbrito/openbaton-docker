package server

import (
	"errors"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

// dummy flavours
var (
	dockerFlavour = &pop.Flavour{
		Id:   "docker-flavour-id",
		Name: "docker.container",
	}

	flavours = map[string]*pop.Flavour{
		dockerFlavour.Id: dockerFlavour,
	}

	flavourNames = map[string]string{
		dockerFlavour.Name: dockerFlavour.Id,
	}
)

// Container fetches created or running containers from Docker, applying the given filter.
func (svc *service) Containers(ctx context.Context, filter *pop.Filter) (*pop.ContainerList, error) {
	tag := util.FuncName()
	op := "Containers"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *filter,
	}).Debug("fetching containers")

	// filter for a container with the given ID
	if filter.Options != nil {
		cont, err := svc.getSingleContainerInfo(ctx, filter)
		if err != nil {
			svc.WithError(err).WithFields(log.Fields{
				"tag":    tag,
				"op":     op,
				"filter": *filter,
			}).Error("fetching single container failed")

			return nil, err
		}

		return &pop.ContainerList{
			List: []*pop.Container{cont},
		}, nil
	}

	return svc.getContainerInfos(ctx)
}

// Flavours are not necessary; the only reason they are implemented it's because they exist in the
// OpenStack/Amazon/... world, and so the NFVO expects one of them.
// We're letting the PoP declare fake flavours, giving an appearance of continuity with the rest of the NFV world.
func (svc *service) Flavours(ctx context.Context, filter *pop.Filter) (*pop.FlavourList, error) {
	tag := util.FuncName()
	op := "Flavours"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *filter,
	}).Debug("fetching flavours")

	if filter.Options != nil {
		fl, err := svc.getSingleFlavourInfo(ctx, filter)
		if err != nil {
			svc.WithError(err).WithFields(log.Fields{
				"tag":    tag,
				"op":     op,
				"filter": *filter,
			}).Error("fetching single flavour failed")

			return nil, err
		}

		return &pop.FlavourList{
			List: []*pop.Flavour{fl},
		}, nil
	}

	return svc.getFlavourInfos(ctx)
}

// Images retrieves and returns the available images on the Docker daemon.
func (svc *service) Images(ctx context.Context, filter *pop.Filter) (*pop.ImageList, error) {
	tag := util.FuncName()
	op := "Images"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *filter,
	}).Debug("fetching images")

	// filter for an image with the given ID or name
	if filter.Options != nil {
		img, err := svc.getSingleImageInfo(ctx, filter)
		if err != nil {
			svc.WithError(err).WithFields(log.Fields{
				"tag":    tag,
				"op":     op,
				"filter": *filter,
			}).Error("fetching single images failed")

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
	tag := util.FuncName()
	op := "Images"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *filter,
	}).Debug("fetching networks")

	// filter for a network with the given ID or name.
	if filter.Options != nil {
		netw, err := svc.getSingleNetworkInfo(ctx, filter)
		if err != nil {
			svc.WithError(err).WithFields(log.Fields{
				"tag":    tag,
				"op":     op,
				"filter": *filter,
			}).Error("fetching single networks failed")

			return nil, err
		}

		return &pop.NetworkList{
			List: []*pop.Network{netw},
		}, nil
	}

	return svc.getNetworkInfos(ctx)
}

// filterContainer uses a Filter to match a container. Because it accesses the cont and names maps, it requires
// at least a read lock to be executed (write lock is also fine for obvious reasons).
func (svc *service) filterContainer(filter *pop.Filter) (*svcCont, error) {
	// id and name can't be both set.
	// If name is set, GetId will return ""
	id := filter.GetId()

	if name := filter.GetName(); name != "" {
		id = svc.names[name] // will set id to "" if the name is not found
	}

	if id == "" {
		return nil, errors.New("no container specified")
	}

	pcont, ok := svc.conts[id]
	if !ok {
		return nil, ErrNoSuchContainer
	}

	return pcont, nil
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

// getFlavourInfos is defined as a method for coherence with the other functions,
// but it's unnecessary.
func (svc *service) getFlavourInfos(context.Context) (*pop.FlavourList, error) {
	list := make([]*pop.Flavour, 0, len(flavours))

	for _, fl := range flavours {
		list = append(list, fl)
	}

	return &pop.FlavourList{List: list}, nil
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

func (svc *service) getSingleContainerInfo(_ context.Context, filter *pop.Filter) (*pop.Container, error) {
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	pcont, err := svc.filterContainer(filter)
	if err != nil {
		return nil, err
	}

	return pcont.Container, nil
}

func (*service) getSingleFlavourInfo(_ context.Context, filter *pop.Filter) (*pop.Flavour, error) {
	id := filter.GetId() // will return "" if not set

	if name := filter.GetName(); name != "" {
		id = flavourNames[name] // will return "" if not found
	}

	fl := flavours[id] // will return nil if not found

	if fl == nil {
		return nil, ErrNoSuchFlavour
	}

	return fl, nil
}

func (svc *service) getSingleImageInfo(ctx context.Context, filter *pop.Filter) (*pop.Image, error) {
	query := filter.GetName()

	if id := filter.GetId(); id != "" {
		query = id // prioritise IDs
	}

	// Docker API in some places accepts either a name or an ID.
	dimg, _, err := svc.cln.ImageInspectWithRaw(ctx, query)
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

func (svc *service) getSingleNetworkInfo(ctx context.Context, filter *pop.Filter) (*pop.Network, error) {
	query := filter.GetName()

	if id := filter.GetId(); id != "" {
		query = id // prioritise IDs
	}

	// Docker API in some places accepts either a name or an ID.
	dnet, err := svc.cln.NetworkInspect(ctx, query)
	if err != nil {
		return nil, err
	}

	return extractNetwork(dnet), nil
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
