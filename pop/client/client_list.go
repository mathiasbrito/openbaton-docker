package client

import (
	"context"
	"fmt"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/openbaton-docker/pop"
)

// Image returns the image on the server having the given id as an OpenBaton NFVImage struct.
func (cln *Client) Image(ctx context.Context, id string) (*catalogue.NFVImage, error) {
	imgs, err := cln.fetchImages(ctx, &pop.Filter{Id: id})
	if err != nil {
		return nil, err
	}

	switch len(imgs) {
	case 0:
		return nil, nil

	case 1:
		return imgs[0], nil

	default:
		return nil, fmt.Errorf("too many images returned from query")
	}
}

// Images returns the images on the server as OpenBaton NFVImage structs.
func (cln *Client) Images(ctx context.Context) ([]*catalogue.NFVImage, error) {
	return cln.fetchImages(ctx, &pop.Filter{})
}

// Server returns the container on the server having the given id as an OpenBaton Server struct.
func (cln *Client) Server(ctx context.Context, id string) (*catalogue.Server, error) {
	srvs, err := cln.fetchServers(ctx, &pop.Filter{Id: id})
	if err != nil {
		return nil, err
	}

	switch len(srvs) {
	case 0:
		return nil, nil

	case 1:
		return srvs[0], nil

	default:
		return nil, fmt.Errorf("too many servers returned from query")
	}
}

// Servers returns the containers on the server as OpenBaton Server structs.
func (cln *Client) Servers(ctx context.Context) ([]*catalogue.Server, error) {
	return cln.fetchServers(ctx, &pop.Filter{})
}

func (cln *Client) fetchImages(ctx context.Context, filter *pop.Filter) ([]*catalogue.NFVImage, error) {
	var imgs []*pop.Image

	op := func(stub pop.PopClient) error {
		ilist, err := stub.Images(ctx, filter)
		if err != nil {
			return err
		}

		if ilist == nil {
			imgs = []*pop.Image{}
		} else {
			imgs = ilist.List
		}

		return nil
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeImages(imgs), nil
}

func (cln *Client) fetchServers(ctx context.Context, filter *pop.Filter) ([]*catalogue.Server, error) {
	var conts []*pop.Container

	op := func(stub pop.PopClient) error {
		clist, err := stub.Containers(ctx, filter)
		if err != nil {
			return err
		}

		if clist == nil {
			conts = []*pop.Container{}
		} else {
			conts = clist.List
		}

		return nil
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeServers(ctx, conts)
}

func (cln *Client) makeImage(img *pop.Image) *catalogue.NFVImage {
	name := ""
	if img.Names != nil && len(img.Names) > 0 {
		name = img.Names[0]
	}

	return &catalogue.NFVImage{
		ExtID:   img.Id,
		Name:    name,
		Created: catalogue.UnixDate(img.Created),
	}
}

func (cln *Client) makeImages(imgs []*pop.Image) []*catalogue.NFVImage {
	nfvImgs := make([]*catalogue.NFVImage, len(imgs))

	for i, img := range imgs {
		nfvImgs[i] = cln.makeImage(img)
	}

	return nfvImgs
}

func (cln *Client) makeServer(ctx context.Context, cont *pop.Container) (srv *catalogue.Server, err error) {
	var nfvImage *catalogue.NFVImage
	if cont.ImageId != "" {
		nfvImage, err = cln.Image(ctx, cont.ImageId)
		if err != nil {
			return nil, err
		}
	}

	name := ""
	if cont.Names != nil && len(cont.Names) > 0 {
		name = cont.Names[0]
	}

	ipMap := make(map[string][]string)

	if cont.Endpoints != nil {
		for netname, ep := range cont.Endpoints {
			if ep.Ipv4 != nil && ep.Ipv4.Address != "" {
				// no ipv6 for now
				ipMap[netname] = []string{ep.Ipv4.Address}
			}
		}
	}

	return &catalogue.Server{
		ExtID:          cont.Id,
		Name:           name,
		Status:         cont.Status,
		ExtendedStatus: cont.ExtendedStatus,
		Image:          nfvImage,
		IPs:            ipMap,
		FloatingIPs:    map[string]string{},
	}, nil
}

func (cln *Client) makeServers(ctx context.Context, conts []*pop.Container) ([]*catalogue.Server, error) {
	servs := make([]*catalogue.Server, len(conts))

	for i, cont := range conts {
		serv, err := cln.makeServer(ctx, cont)
		if err != nil {
			return nil, err
		}

		servs[i] = serv
	}

	return servs, nil
}
