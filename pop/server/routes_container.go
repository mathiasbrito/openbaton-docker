package server

import (
    "errors"
    "fmt"
    "time"

    "golang.org/x/net/context"

    pop "github.com/mcilloni/openbaton-docker/pop/proto"
    "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
    "github.com/docker/docker/api/types"
    "github.com/golang/protobuf/ptypes/empty"
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
                NetworkID: endp.NetId,
                EndpointID: endp.EndpointId,
                IPAMConfig: ipcfg,
            }
        }
    }

    ccb, err := svc.cln.ContainerCreate(
        ctx, 
        &container.Config{
            Hostname: cfg.Name,
            Image: cfg.ImageId,
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
