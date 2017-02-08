package server

import (
	"bytes"
	"errors"
	"fmt"
	"net"
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
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

var (
	ErrAlreadyStarted  = grpc.Errorf(codes.AlreadyExists, "container already started")
	ErrAlreadyStopped  = grpc.Errorf(codes.Unavailable, "container already stopped")
	ErrInvalidArgument = grpc.Errorf(codes.InvalidArgument, "invalid argument provided")
	ErrInvalidState    = grpc.Errorf(codes.FailedPrecondition, "container is in an invalid state")
	ErrNoSuchContainer = grpc.Errorf(codes.NotFound, "no container for the given filter")
	ErrNoSuchFlavour   = grpc.Errorf(codes.NotFound, "no flavour for the given filter")
	ErrNoSuchNetwork   = grpc.Errorf(codes.NotFound, "no network for the given filter")
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
	tag := util.FuncName()
	op := "Create"

	svc.WithFields(log.Fields{
		"tag":  tag,
		"op":   op,
		"name": cfg.Name,
	}).Debug("attempting to create container")

	if cfg.FlavourId != "" && cfg.FlavourId != dockerFlavour.Id {
		svc.WithFields(log.Fields{
			"tag":  tag,
			"op":   op,
			"name": cfg.Name,
		}).Error("cannot create container")

		return nil, fmt.Errorf("unsupported flavour %v, only %v is supported", cfg.FlavourId, dockerFlavour.Id)
	}

	// grab the lock BEFORE creating the descriptor -
	// if cfg specifies an already taken name, it's pointless to
	// waste time with the daemon. To check this we need the lock to read from names.
	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	svc.WithFields(log.Fields{
		"tag":  tag,
		"op":   op,
		"name": cfg.Name,
	}).Debug("list lock obtained")

	if _, found := svc.contNames[cfg.Name]; found {
		return nil, grpc.Errorf(codes.AlreadyExists, "container name %s already taken", cfg.Name)
	}

	cont, err := svc.createPcont(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// after creating a new container descriptor, find an id and
	// put it into the map.

	// check if the name is al

	for {
		cont.Id = util.GenerateID()

		if _, found := svc.conts[cont.Id]; !found {
			break
		}
	}

	// add the container to the container list of the server.
	svc.conts[cont.Id] = &svcCont{
		Container: cont,
		DockerID:  "", // not yet assigned
	}

	// cont.Names has always at least an element
	svc.contNames[cont.Names[0]] = cont.Id

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-id":    cont.Id,
		"cont-names": cont.Names,
	}).Debug("created container")

	return cont, nil
}

// Delete removes the containers identified by the given filter, stopping it before if necessary.
func (svc *service) Delete(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	tag := util.FuncName()
	op := "Delete"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("attempting to delete a container")

	// get the lock before editing the map
	svc.contsMux.Lock()
	defer svc.contsMux.Unlock()

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("list lock obtained")

	pcont, err := svc.filterContainer(filter)
	if err != nil {
		svc.WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": filter.Options != nil,
		}).Error("container fetching failed")

		return nil, err
	}

	// The sync.Mutex avoids race conditions while deleting the container.
	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	if pcont.Status == pop.Container_RUNNING {
		svc.WithFields(log.Fields{
			"tag":  tag,
			"op":   op,
			"cont": pcont.Names,
		}).Debug("trying to stop container")

		if err := svc.stopContainer(ctx, pcont); err != nil && err != ErrAlreadyStopped {
			svc.WithFields(log.Fields{
				"tag":  tag,
				"op":   op,
				"cont": pcont.Names,
			}).Error("container stopping failed")
			return nil, err
		}

		svc.WithFields(log.Fields{
			"tag":  tag,
			"op":   op,
			"cont": pcont.Names,
		}).Debug("container stopped")
	}

	// deletes the container from the container list and from names
	delete(svc.conts, pcont.Id)
	delete(svc.contNames, pcont.Names[0])

	// if someone still holds a reference to this container
	pcont.Status = pop.Container_UNAVAILABLE
	pcont.ExtendedStatus = "this container has been deleted"

	svc.WithFields(log.Fields{
		"tag":  tag,
		"op":   op,
		"cont": pcont.Names,
	}).Debug("trying to stop container")

	return &empty.Empty{}, nil
}

// Metadata adds the given metadata values to the container that matches with the ID.
// An empty value for a key means that the key will be removed from the metadata.
// Metadata will return an error if the container has already been spawned.
func (svc *service) Metadata(ctx context.Context, newMD *pop.NewMetadata) (*empty.Empty, error) {
	tag := util.FuncName()
	op := "Metadata"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *newMD.Filter,
	}).Debug("attempting to add metadata to a container")

	if newMD.Filter == nil {
		return nil, errors.New("empty filter")
	}

	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": *newMD.Filter,
	}).Debug("list rlock obtained")

	pcont, err := svc.filterContainer(newMD.Filter)
	if err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": *newMD.Filter,
		}).Error("cannot retrieve container")

		return nil, err
	}

	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-names": pcont.Names,
	}).Debug("container lock obtained")

	if pcont.Status != pop.Container_CREATED {
		svc.WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Error("cannot set metadata, container has already started")

		return nil, ErrAlreadyStarted
	}

	if newMD.Md == nil || newMD.Md.Entries == nil {
		svc.WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Debug("container lock obtained")
		return nil, ErrInvalidArgument
	}

	pcont.Md().Merge(newMD.Md.Entries)

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-names": pcont.Names,
	}).Debug("container lock obtained")

	return &empty.Empty{}, nil
}

// Start starts the container identified by the given filter, by creating and launching a Docker
// container with its metadata as environment variables.
func (svc *service) Start(ctx context.Context, filter *pop.Filter) (*pop.Container, error) {
	tag := util.FuncName()
	op := "Start"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("attempting to start a container")

	// In case a container is quickly created and then started, this avoids races.
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("list lock obtained")

	pcont, err := svc.filterContainer(filter)
	if err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": filter.Options != nil,
		}).Error("cannot find container")

		return nil, err
	}

	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-names": pcont.Names,
	}).Debug("container lock obtained")

	// Ensures the container is launched once and only once.
	switch pcont.Status {
	case pop.Container_EXITED:
		fallthrough
	case pop.Container_RUNNING:
		svc.WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Error("cannot start a container already started")

		return nil, ErrAlreadyStarted

	case pop.Container_CREATED:
		// go through

	default:
		svc.WithFields(log.Fields{
			"tag":         tag,
			"op":          op,
			"cont-names":  pcont.Names,
			"cont-status": pcont.Status.String(),
		}).Error("invalid status")

		return nil, ErrInvalidState
	}

	ccb, err := svc.createContainer(ctx, pcont)
	if err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Error("cannot create container")

		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while creating: %v", err)
		return nil, err
	}

	if len(ccb.Warnings) != 0 {
		pcont.ExtendedStatus = fmt.Sprintf("warnings from container instantiation: [%s]", strings.Join(ccb.Warnings, ", "))
	}

	if err := svc.cln.ContainerStart(ctx, ccb.ID, types.ContainerStartOptions{}); err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Error("cannot start container")

		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while starting: %v", err)
		return nil, err
	}

	if err := svc.updateContainer(ctx, pcont, ccb.ID); err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":        tag,
			"op":         op,
			"cont-names": pcont.Names,
		}).Warn("cannot update container")

		pcont.ExtendedStatus = "warning: update of this container failed"
		return nil, err
	}

	pcont.Status = pop.Container_RUNNING
	pcont.ExtendedStatus = "the container is running"

	return pcont.ToPop(), nil
}

// Stop stops the container identified by the given filter.
func (svc *service) Stop(ctx context.Context, filter *pop.Filter) (*empty.Empty, error) {
	tag := util.FuncName()
	op := "Stop"

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("attempting to stop a container")

	// In case a container is quickly created and then stopped, this avoids races.
	svc.contsMux.RLock()
	defer svc.contsMux.RUnlock()

	svc.WithFields(log.Fields{
		"tag":    tag,
		"op":     op,
		"filter": filter.Options != nil,
	}).Debug("list rlock obtained")

	pcont, err := svc.filterContainer(filter)
	if err != nil {
		svc.WithError(err).WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": filter.Options != nil,
		}).Error("cannot find container")

		return nil, err
	}

	// Get a lock on this container, to safely handle its state
	pcont.mux.Lock()
	defer pcont.mux.Unlock()

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-names": pcont.Names,
	}).Debug("list rlock obtained")

	switch pcont.Status {
	case pop.Container_EXITED:
		fallthrough
	case pop.Container_FAILED:
		svc.WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": filter.Options != nil,
		}).Error("cannot stop a stopped container")

		return nil, ErrAlreadyStopped

	case pop.Container_RUNNING:
		// go through

	case pop.Container_CREATED:
		svc.WithFields(log.Fields{
			"tag":    tag,
			"op":     op,
			"filter": filter.Options != nil,
		}).Error("container was never started")

		return nil, ErrNotStarted

	default:
		svc.WithFields(log.Fields{
			"tag":         tag,
			"op":          op,
			"cont-names":  pcont.Names,
			"cont-status": pcont.Status.String(),
		}).Error("invalid status")

		return nil, ErrInvalidState
	}

	// The switch above ensures the container is stopped once and only once.

	svc.WithFields(log.Fields{
		"tag":        tag,
		"op":         op,
		"cont-names": pcont.Names,
	}).Debug("container stopped")

	return &empty.Empty{}, svc.stopContainer(ctx, pcont)
}

func (svc *service) createContainer(ctx context.Context, pcont *svcCont) (ccb container.ContainerCreateCreatedBody, opErr error) {
	var dockerEndpoints map[string]*network.EndpointSettings

	if pcont.Endpoints != nil {
		dockerEndpoints = make(map[string]*network.EndpointSettings)

		for netID, endp := range pcont.Endpoints {
			dockerEndpoints[netID] = &network.EndpointSettings{
				NetworkID:  endp.NetId,
				EndpointID: endp.EndpointId,
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: endp.Ipv4.Address, // ensured by creatPcont to be valid
				},
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
			Env:      pcont.Md().Strings(),
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

func (svc *service) createPcont(ctx context.Context, cfg *pop.ContainerConfig) (*pop.Container, error) {
	if cfg.ImageId == "" {
		return nil, errors.New("no image ID provided")
	}

	// check if the image exists
	filter := &pop.Filter{Options: &pop.Filter_Id{Id: cfg.ImageId}}

	if _, err := svc.getSingleImageInfo(ctx, filter); err != nil {
		return nil, err
	}

	endpoints := cfg.Endpoints
	if endpoints == nil || len(endpoints) == 0 {
		// associate to the default network
		ep, err := svc.getDefaultEndpoint()
		if err != nil {
			return nil, err
		}

		endpoints = map[string]*pop.Endpoint{
			svc.defaultNet.ID: ep,
		}
	} else {
		// allocate IPs on networks
		for netID, endp := range endpoints {
			pnet, found := svc.nets[netID]
			if !found {
				return nil, grpc.Errorf(codes.NotFound, "no network with id %s", netID)
			}

			// update endpoint network name
			endp.NetName = pnet.Name

			// no IPv6 support atm

			if endp.Ipv4 != nil && endp.Ipv4.Address != "" {
				// Reserve the ipv4
				if err := pnet.ReserveV4(net.ParseIP(endp.Ipv4.Address)); err != nil {
					return nil, err
				}
			} else {
				// dynamically allocate a new IP

				ip, _, err := pnet.GetV4()
				if err != nil {
					return nil, err
				}

				endp.Ipv4 = &pop.Ip{
					Address: ip.String(),
					Subnet:  pnet.ToPopSubnet(),
				}
			}
		}
	}

	return &pop.Container{
		Status:         pop.Container_CREATED,
		ExtendedStatus: "container ready for instantiation",
		Names:          []string{cfg.Name},
		ImageId:        cfg.ImageId,
		FlavourId:      cfg.FlavourId,
		Created:        time.Now().Unix(),
		Endpoints:      endpoints,
		Md:             &pop.Metadata{Entries: make(map[string]string)},
	}, nil
}

// stopContainer stops a container; this function expects to hold the lock on the given pcont.
func (svc *service) stopContainer(ctx context.Context, pcont *svcCont) error {
	timeout := 5 * time.Second

	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
	}

	if err := svc.cln.ContainerStop(ctx, pcont.DockerID, &timeout); err != nil {
		pcont.Status = pop.Container_FAILED
		pcont.ExtendedStatus = fmt.Sprintf("error while stopping: %v", err)
		return err
	}

	pcont.Status = pop.Container_EXITED
	pcont.ExtendedStatus = "the container has exited"
	pcont.DockerID = ""

	// in the end, release the container IPs
	return svc.releaseContIPs(pcont)
}

func (svc *service) updateContainer(ctx context.Context, pcont *svcCont, dockerID string) error {
	dcont, err := svc.cln.ContainerInspect(ctx, dockerID)
	if err != nil {
		return err
	}

	started, err := time.Parse(time.RFC3339Nano, dcont.Created)
	if err != nil {
		return pop.InternalErr
	}

	b := bytes.Buffer{}
	for _, part := range dcont.Config.Cmd {
		b.WriteString(part)
		b.WriteRune(' ')
	}

	pcont.DockerID = dockerID
	pcont.Command = b.String()
	pcont.FlavourId = dockerFlavour.Id
	pcont.Started = started.Unix()

	pcont.Endpoints = extractEndpoints(dcont.NetworkSettings.Networks)

	return nil
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
