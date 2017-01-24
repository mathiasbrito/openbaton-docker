package client

import (
    "context"
    "errors"

    pop "github.com/mcilloni/openbaton-docker/pop/proto"
    "github.com/mcilloni/go-openbaton/catalogue"
)

// Create creates a new server on the remote Pop. An entry in "ips" with an entry IP will randomly assign an IP from the given network.
func (cln *Client) Create(ctx context.Context, name, imageID, flavorID string, ips map[string]string) (*catalogue.Server, error) {
	var cont *pop.Container

    endpoints := make(map[string]*pop.Endpoint)
    for net, ip := range ips {
        endpoints[net] = &pop.Endpoint {
            Ipv4: &pop.Ip{Address: ip},
        }
    }

    cfg := &pop.ContainerConfig{
        Name: name,
        ImageId: imageID,
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

// Start starts a Server created by Create().
func (cln *Client) Start(ctx context.Context, id string) (*catalogue.Server, error) {
    var cont *pop.Container
    
    op := func(stub pop.PopClient) (err error) {
        cont, err = stub.Start(ctx, &pop.Filter{Id: id})
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
