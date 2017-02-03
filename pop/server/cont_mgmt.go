package server

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

    pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/container"
	"github.com/satori/go.uuid"
	"github.com/docker/docker/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrAlreadyLaunched = grpc.Errorf(codes.AlreadyExists, "container already spawned")
	ErrInvalidArgument = grpc.Errorf(codes.InvalidArgument, "invalid argument provided")
	ErrNoSuchContainer = grpc.Errorf(codes.NotFound, "no container for the given ID")
)

func (svc *service) containers() []*pop.Container {
    svc.contsMux.RLock()
    defer svc.contsMux.RUnlock()

    // preallocate the list
    conts := make([]*pop.Container, 0,  len(svc.conts))

    for _, cj := range svc.conts {
        conts = append(conts, cj.Container)
    }

    return conts
}

func (svc *service) createContainer(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.FlavourId != "" && cfg.FlavourId != dockerFlavour.Id {
		return nil, fmt.Errorf("unsupported flavour %v, only %v is supported", cfg.FlavourId, dockerFlavour.Id)
	}		

	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	if cfg.ImageId == "" {
		return nil, errors.New("no image ID provided")
	}

	// check if the image exists
	if _, err := svc.getSingleImageInfo(ctx, cfg.ImageId); err != nil {
		return nil, err
	}

	var id string
	for {
		id = uuid.NewV4().String()
		
		if _, found := svc.conts[id]; !found {
			break
		}
	}

	cont := &pop.Container{
		Id: id,
		Status: pop.Container_CREATED,
		ExtendedStatus: "container ready for instantiation",
		Names: []string{cfg.Name},
		ImageId: cfg.ImageId,
		FlavourId: cfg.FlavourId,
		Created: time.Now().Unix(),
		Endpoints: cfg.Endpoints,
	}
	
	// add the container to the container list of the server.
	svc.conts[id] = &svcCont{
		Container: cont,
		DockerID: "", // not yet assigned

		metadata: make(metadata),
	} 

	return cont, nil
}

func (svc *service) launchDockerContainer(ctx context.Context, id string) (cont *pop.Container, actionErr error) {
	pcont, ok := svc.conts[id]
	if !ok {
		return nil, fmt.Errorf("launchDockerContainer: no container with id %s", id)
	}

	actionErr = ErrAlreadyLaunched

	// Ensures the container is launched once and only once. 
	// The closure will set the correct return values in case it is called;
	// otherwise, ErrAlreadyLaunched will be returned.
	pcont.launch.Do(func() {
		var dockerEndpoints map[string]*network.EndpointSettings

		if pcont.Endpoints != nil {
			dockerEndpoints = make(map[string]*network.EndpointSettings)

			for netname, endp := range pcont.Endpoints {
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
				// Names should have at least one member, otherwise 
				// there is a bug
				Hostname: pcont.Names[0],
				Image:    pcont.ImageId,
				Env: pcont.metadata.Strings(),
			},
			&container.HostConfig{
				AutoRemove: true,
			},
			&network.NetworkingConfig{
				EndpointsConfig: dockerEndpoints,
			},
			pcont.Names[0],
		)

		if err != nil {
			actionErr = err
			pcont.Status = pop.Container_DEAD
			pcont.ExtendedStatus = fmt.Sprintf("error while creating: %v", err)
			return 
		}

		pcont.DockerID = ccb.ID
		if len(ccb.Warnings) != 0 {
			pcont.ExtendedStatus = fmt.Sprintf("warnings from container instantiation: [%s]", strings.Join(ccb.Warnings, ", "))
		}

		if err := svc.cln.ContainerStart(ctx, ccb.ID, types.ContainerStartOptions{}); err != nil {
			actionErr = err
			pcont.Status = pop.Container_DEAD
			pcont.ExtendedStatus = fmt.Sprintf("error while starting: %v", err)
			return
		}

		cont = pcont.Container
		actionErr = nil
	})

	return
}