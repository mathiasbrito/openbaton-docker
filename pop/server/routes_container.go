package server

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/golang/protobuf/ptypes/empty"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

// Create creates a new container as described by the given config.
func (svc *service) Create(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.FlavourId != "" && cfg.FlavourId != dockerFlavour.Id {
		return nil, fmt.Errorf("unsupported flavour %v, only %v is supported", cfg.FlavourId, dockerFlavour.Id)
	}

	var dockerEndpoints map[string]*network.EndpointSettings

	if cfg.Endpoints != nil {
		dockerEndpoints = make(map[string]*network.EndpointSettings)

		for netname, endp := range cfg.Endpoints {
			var ipcfg *network.EndpointIPAMConfig

			if endp.Ipv4.Address != "" || endp.Ipv6.Address != "" {
				ipcfg = &network.EndpointIPAMConfig{
					IPv4Address: endp.Ipv4.Address,
					IPv6Address: endp.Ipv6.Address,
				}
			}

			dockerEndpoints[netname] = &network.EndpointSettings{
				NetworkID:  endp.NetId,
				EndpointID: endp.EndpointId,
				IPAMConfig: ipcfg,
			}
		}
	}

	ccb, err := svc.cln.ContainerCreate(
		ctx,
		&container.Config{
			Hostname: cfg.Name,
			Image:    cfg.ImageId,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		&network.NetworkingConfig{
			EndpointsConfig: dockerEndpoints,
		},
		cfg.Name,
	)

	if err != nil {
		return nil, err
	}

	return svc.getSingleContainerInfo(ctx, ccb.ID)
}

// Delete removes the containers identified by the given filter.
func (svc *service) Delete(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	if filter.Id == "" {
		return nil, errors.New("no container specified for Start")
	}

	timeout := new(time.Duration)
	*timeout = time.Minute

	if deadline, ok := ctx.Deadline(); ok {
		*timeout = deadline.Sub(time.Now())
	}

	if err := svc.cln.ContainerStop(ctx, filter.Id, timeout); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Metadata adds the given metadata values to the container that matches with the ID.
// An empty value for a key means that the key will be removed from the metadata.
// Metadata will return an error if the container has already been spawned.
func (svc *service) Metadata(ctx context.Context, newMD *pop.NewMetadata) (*empty.Empty, error) {
	pcont := svc.conts[newMD.Id]
	if pcont == nil {
		return nil, ErrNoSuchContainer
	}
	
	if pcont.Status != pop.Container_CREATED {
		return nil, ErrAlreadyLaunched
	}

	if newMD.Md == nil || newMD.Md.Entries == nil {
		return nil, ErrInvalidArgument
	}

	pcont.metadata.Merge(newMD.Md.Entries)

	return &empty.Empty{}, nil
}

// Start starts the container identified by the given filter.
func (svc *service) Start(ctx context.Context, filter *pop.Filter) (*pop.Container, error) {
	if filter.Id == "" {
		return nil, errors.New("no container specified for Start")
	}

	if err := svc.cln.ContainerStart(ctx, filter.Id, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return svc.getSingleContainerInfo(ctx, filter.Id)
}
