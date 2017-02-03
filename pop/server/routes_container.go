package server

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/golang/protobuf/ptypes/empty"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/satori/go.uuid"
)

var (
	ErrAlreadyStarted  = grpc.Errorf(codes.AlreadyExists, "container already started")
	ErrAlreadyStopped  = grpc.Errorf(codes.Unavailable, "container already stopped")
	ErrInvalidArgument = grpc.Errorf(codes.InvalidArgument, "invalid argument provided")
	ErrInvalidState    = grpc.Errorf(codes.FailedPrecondition, "container is in an invalid state")
	ErrNoSuchContainer = grpc.Errorf(codes.NotFound, "no container for the given ID")
	ErrNotStarted      = grpc.Errorf(codes.Unavailable, "container not started yet")
)

// A container can have one of the following strictly sequential life cycles:
// Created -> Exited
// Created -> Failed
// Created -> Running -> Failed
// Created -> Running -> Failed
// Created -> Running -> Exited
//
// The succession of these states is enforced through a per-container mutex.

// Create creates a new container as described by the given config.
func (svc *service) Create(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.FlavourId != "" && cfg.FlavourId != dockerFlavour.Id {
		return nil, fmt.Errorf("unsupported flavour %v, only %v is supported", cfg.FlavourId, dockerFlavour.Id)
	}

	cont, err := svc.checkConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// after creating a new container descriptor, get the lock, find an id and
	// put it into the map.
	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	for {
		cont.Id = uuid.NewV4().String()

		if _, found := svc.conts[cont.Id]; !found {
			break
		}
	}

	// add the container to the container list of the server.
	svc.conts[cont.Id] = &svcCont{
		Container: cont,
		DockerID:  "", // not yet assigned

		metadata: make(metadata),
	}

	return cont, nil
}

// Delete removes the containers identified by the given filter, stopping it before if necessary.
func (svc *service) Delete(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	if filter.Id == "" {
		return nil, errors.New("no container specified for Delete")
	}

	// get the lock before editing the map
	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	pcont, ok := svc.conts[filter.Id]
	if !ok {
		return nil, ErrNoSuchContainer
	}

	// The sync.Mutex avoids race conditions while deleting the container.
	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	if pcont.Status == pop.Container_RUNNING {
		if err := svc.stopContainer(ctx, pcont); err != nil && err != ErrAlreadyStopped {
			return nil, err
		}
	}

	// deletes the container from the container list
	delete(svc.conts, filter.Id)

	// if someone still holds a reference to this container
	pcont.Status = pop.Container_EXITED
	return &empty.Empty{}, nil
}

// Metadata adds the given metadata values to the container that matches with the ID.
// An empty value for a key means that the key will be removed from the metadata.
// Metadata will return an error if the container has already been spawned.
func (svc *service) Metadata(ctx context.Context, newMD *pop.NewMetadata) (*empty.Empty, error) {
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	pcont := svc.conts[newMD.Id]
	if pcont == nil {
		return nil, ErrNoSuchContainer
	}

	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	if pcont.Status != pop.Container_CREATED {
		return nil, ErrAlreadyStarted
	}

	if newMD.Md == nil || newMD.Md.Entries == nil {
		return nil, ErrInvalidArgument
	}

	pcont.metadata.Merge(newMD.Md.Entries)

	return &empty.Empty{}, nil
}

// Start starts the container identified by the given filter, by creating and launching a Docker
// container with its metadata as environment variables.
func (svc *service) Start(ctx context.Context, filter *pop.Filter) (*pop.Container, error) {
	if filter.Id == "" {
		return nil, errors.New("no container specified for Start")
	}

	// In case a container is quickly created and then started, this avoids races.
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	pcont, ok := svc.conts[filter.Id]
	if !ok {
		return nil, fmt.Errorf("start: no container with id %s", filter.Id)
	}

	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	// Ensures the container is launched once and only once.
	switch pcont.Status {
	case pop.Container_EXITED:
		fallthrough
	case pop.Container_RUNNING:
		return nil, ErrAlreadyStarted

	case pop.Container_CREATED:
		// go through

	default:
		return nil, ErrInvalidState
	}

	ccb, err := svc.createContainer(ctx, pcont)
	if err != nil {
		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while creating: %v", err)
		return nil, err
	}

	pcont.DockerID = ccb.ID
	if len(ccb.Warnings) != 0 {
		pcont.ExtendedStatus = fmt.Sprintf("warnings from container instantiation: [%s]", strings.Join(ccb.Warnings, ", "))
	}

	if err := svc.cln.ContainerStart(ctx, ccb.ID, types.ContainerStartOptions{}); err != nil {
		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while starting: %v", err)
		return nil, err
	}

	pcont.Status = pop.Container_RUNNING

	return pcont.Container, nil
}

// Stop stops the container identified by the given filter.
func (svc *service) Stop(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	if filter.Id == "" {
		return nil, errors.New("no container specified for Start")
	}

	// In case a container is quickly created and then stopped, this avoids races.
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	pcont, ok := svc.conts[filter.Id]
	if !ok {
		return nil, fmt.Errorf("stop: no container with id %s", filter.Id)
	}

	// Get a lock on this container, to safely handle its state
	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	switch pcont.Status {
	case pop.Container_EXITED:
		fallthrough
	case pop.Container_FAILED:
		return nil, ErrAlreadyStopped

	case pop.Container_RUNNING:
		// go through

	case pop.Container_CREATED:
		return nil, ErrNotStarted

	default:
		return nil, ErrInvalidState
	}

	// The switch above ensures the container is stopped once and only once.

	return &empty.Empty{}, svc.stopContainer(ctx, pcont)
}

func (svc *service) checkConfig(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.ImageId == "" {
		return nil, errors.New("no image ID provided")
	}

	// check if the image exists
	if _, err := svc.getSingleImageInfo(ctx, cfg.ImageId); err != nil {
		return nil, err
	}

	return &pop.Container{
		Status:         pop.Container_CREATED,
		ExtendedStatus: "container ready for instantiation",
		Names:          []string{cfg.Name},
		ImageId:        cfg.ImageId,
		FlavourId:      cfg.FlavourId,
		Created:        time.Now().Unix(),
		Endpoints:      cfg.Endpoints,
	}, nil
}

func (svc *service) createContainer(ctx context.Context, pcont *svcCont) (container.ContainerCreateCreatedBody, error) {
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

	return svc.cln.ContainerCreate(
		ctx,
		&container.Config{
			// Names should have at least one member, otherwise
			// there is a bug
			Hostname: pcont.Names[0],
			Image:    pcont.ImageId,
			Env:      pcont.metadata.Strings(),
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		&network.NetworkingConfig{
			EndpointsConfig: dockerEndpoints,
		},
		pcont.Names[0],
	)
}

// stopContainer stops a container; this function expects to hold the lock on the given pcont.
func (svc *service) stopContainer(ctx context.Context, pcont *svcCont) error {
	timeout := time.Minute

	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
	}

	if err := svc.cln.ContainerStop(ctx, pcont.DockerID, &timeout); err != nil {
		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while stopping: %v", err)
		return err
	}

	pcont.Status = pop.Container_EXITED

	return nil
}
