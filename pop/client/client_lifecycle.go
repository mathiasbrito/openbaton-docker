package client

import (
	"context"
	"errors"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/catalogue"
)

// AddMetadata adds metadata keys to a not yet started container.
// New keys will match existing ones; any empty value will delete its key from the map server side.
func (cln *Client) AddMetadata(ctx context.Context, f Filter, keys map[string]string) error {
	op := func(stub pop.PopClient) (err error) {
		_, err = stub.Metadata(
			ctx,
			&pop.NewMetadata{
				Filter: filter(f),
				Md: &pop.Metadata{
					Entries: keys,
				},
			},
		)

		if err != nil {
			return
		}

		return
	}

	if err := cln.doRetry(op); err != nil {
		return err
	}

	return nil
}

// Create creates a new server on the remote Pop. An entry in "ips" with an entry IP will randomly assign an IP from the given network.
func (cln *Client) Create(ctx context.Context, name, imageID, flavorID string, ips map[string]string) (*catalogue.Server, error) {
	var cont *pop.Container

	var endpoints map[string]*pop.Endpoint
	if ips != nil {
		endpoints = make(map[string]*pop.Endpoint)
		for net, ip := range ips {
			endpoints[net] = &pop.Endpoint{
				Ipv4: &pop.Ip{Address: ip},
			}
		}
	}

	cfg := &pop.ContainerConfig{
		Name:      name,
		ImageId:   imageID,
		FlavourId: flavorID,
		Endpoints: endpoints,
	}

	op := func(stub pop.PopClient) (err error) {
		cont, err = stub.Create(ctx, cfg)
		if err != nil {
			return
		}

		if cont == nil {
			return errors.New("no container has been created")
		}

		return
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeServer(ctx, cont)
}

// Delete stops and deletes the container identified by the given filter.
func (cln *Client) Delete(ctx context.Context, f Filter) error {
	op := func(stub pop.PopClient) (err error) {
		_, err = stub.Delete(ctx, filter(f))
		if err != nil {
			return
		}

		return
	}

	if err := cln.doRetry(op); err != nil {
		return err
	}

	return nil
}

// Spawn creates and starts a new server on the remote Pop. An entry in "ips" with an entry IP will randomly assign an IP from the given network.
func (cln *Client) Spawn(ctx context.Context, name, imageID, flavorID string, ips map[string]string) (*catalogue.Server, error) {
	srv, err := cln.Create(ctx, name, imageID, flavorID, ips)
	if err != nil {
		return nil, err
	}

	return cln.Start(ctx, IDFilter(srv.ExtID))
}

// Start starts a Server created by Create().
func (cln *Client) Start(ctx context.Context, f Filter) (*catalogue.Server, error) {
	var cont *pop.Container

	op := func(stub pop.PopClient) (err error) {
		cont, err = stub.Start(ctx, filter(f))
		if err != nil {
			return
		}

		if cont == nil {
			return errors.New("no container has been started")
		}

		return
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeServer(ctx, cont)
}

// Stop stops a Server launched by Start or Spawn.
func (cln *Client) Stop(ctx context.Context, f Filter) error {
	op := func(stub pop.PopClient) (err error) {
		_, err = stub.Stop(ctx, filter(f))
		return
	}

	if err := cln.doRetry(op); err != nil {
		return err
	}

	return nil
}
