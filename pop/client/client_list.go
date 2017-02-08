package client

import (
	"context"
	"fmt"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/catalogue"
)

// Flavour returns the flavour having the given filter as an OpenBaton DeploymentFlavour.
func (cln *Client) Flavour(ctx context.Context, f Filter) (*catalogue.DeploymentFlavour, error) {
	flavs, err := cln.fetchFlavours(ctx, f)
	if err != nil {
		return nil, err
	}

	switch len(flavs) {
	case 0:
		return nil, nil

	case 1:
		return flavs[0], nil

	default:
		return nil, fmt.Errorf("too many flavours returned from query")
	}
}

// Flavours returns all the available flavours as OpenBaton DeploymentFlavour.
func (cln *Client) Flavours(ctx context.Context) ([]*catalogue.DeploymentFlavour, error) {
	return cln.fetchFlavours(ctx, nil)
}

// Image returns the image on the server having the given filter as an OpenBaton NFVImage struct.
func (cln *Client) Image(ctx context.Context, f Filter) (*catalogue.NFVImage, error) {
	imgs, err := cln.fetchImages(ctx, f)
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
	return cln.fetchImages(ctx, nil)
}

// Network returns the network on the server having the given filter as an OpenBaton Network struct.
func (cln *Client) Network(ctx context.Context, f Filter) (*catalogue.Network, error) {
	nets, err := cln.fetchNetworks(ctx, f)
	if err != nil {
		return nil, err
	}

	switch len(nets) {
	case 0:
		return nil, nil

	case 1:
		return nets[0], nil

	default:
		return nil, fmt.Errorf("too many networks returned from query")
	}
}

// Networks returns the networks on the server as OpenBaton Network structs.
func (cln *Client) Networks(ctx context.Context) ([]*catalogue.Network, error) {
	return cln.fetchNetworks(ctx, nil)
}

// Server returns the container on the server having the given id as an OpenBaton Server struct.
func (cln *Client) Server(ctx context.Context, f Filter) (*catalogue.Server, error) {
	srvs, err := cln.fetchServers(ctx, f)
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
	return cln.fetchServers(ctx, nil)
}

// fetchFlavours fetches and converts pop Flavours into DeploymentFlavours.
func (cln *Client) fetchFlavours(ctx context.Context, f Filter) ([]*catalogue.DeploymentFlavour, error) {
	var rflavs []*pop.Flavour

	op := func(stub pop.PopClient) error {
		flist, err := stub.Flavours(ctx, filter(f))
		if err != nil {
			return err
		}

		if flist == nil {
			rflavs = []*pop.Flavour{}
		} else {
			rflavs = flist.List
		}

		return nil
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeFlavours(rflavs), nil
}

// fetchImages fetches and converts pop Images into NFVImages.
func (cln *Client) fetchImages(ctx context.Context, f Filter) ([]*catalogue.NFVImage, error) {
	var imgs []*pop.Image

	op := func(stub pop.PopClient) error {
		ilist, err := stub.Images(ctx, filter(f))
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

	// if there is no such image, end here
	if len(imgs) == 0 {
		return []*catalogue.NFVImage{}, nil
	}

	nfvImgs := cln.makeImages(imgs)

	// if we are filtering for a name, than
	// we need to get just a single image and set the right name to it.
	// if we are filtering for an ID, then just return one of them, they just have
	// different names.
	switch rf := f.(type) {
	case NameFilter:
		nfvImgs[0].Name = string(rf)
		return nfvImgs[:1], nil

	case IDFilter:
		return nfvImgs[:1], nil

	default:
		return nfvImgs, nil
	}
}

// fetchNetworks fetches and converts pop Networks into catalogue.Network instances.
func (cln *Client) fetchNetworks(ctx context.Context, f Filter) ([]*catalogue.Network, error) {
	var rnets []*pop.Network

	op := func(stub pop.PopClient) error {
		nlist, err := stub.Networks(ctx, filter(f))
		if err != nil {
			return err
		}

		if nlist == nil {
			rnets = []*pop.Network{}
		} else {
			rnets = nlist.List
		}

		return nil
	}

	if err := cln.doRetry(op); err != nil {
		return nil, err
	}

	return cln.makeNetworks(rnets), nil
}

// fetchServers gets and creates catalogue.Server instances from pop containers.
func (cln *Client) fetchServers(ctx context.Context, f Filter) ([]*catalogue.Server, error) {
	var conts []*pop.Container

	op := func(stub pop.PopClient) error {
		clist, err := stub.Containers(ctx, filter(f))
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

// makeFlavour converts a pop Flavour into a DeploymentFlavour.
func (cln *Client) makeFlavour(flav *pop.Flavour) *catalogue.DeploymentFlavour {
	return &catalogue.DeploymentFlavour{
		ExtID:      flav.Id,
		FlavourKey: flav.Name,
	}
}

// makeFlavours converts a list of pop Flavour into a list of DeploymentFlavour.
func (cln *Client) makeFlavours(flavs []*pop.Flavour) []*catalogue.DeploymentFlavour {
	depFlavs := make([]*catalogue.DeploymentFlavour, len(flavs))

	for i, flav := range flavs {
		depFlavs[i] = cln.makeFlavour(flav)
	}

	return depFlavs
}

// makeImage converts a pop Image into one (or more) NFVImage.
func (cln *Client) makeImage(img *pop.Image) []*catalogue.NFVImage {
	base := catalogue.NFVImage{
		ExtID:   img.Id,
		Created: catalogue.UnixDate(img.Created),
	}

	if img.Names == nil || len(img.Names) == 0 {
		return []*catalogue.NFVImage{&base}
	}

	ret := make([]*catalogue.NFVImage, 0, len(img.Names))

	// create an image for each tag.
	for _, name := range img.Names {
		img := new(catalogue.NFVImage)
		*img = base

		img.Name = name

		ret = append(ret, img)
	}

	return ret
}

// makeImages converts a list of pop Image into a list of NFVImage.
func (cln *Client) makeImages(imgs []*pop.Image) []*catalogue.NFVImage {
	nfvImgs := make([]*catalogue.NFVImage, 0, len(imgs)) // pre allocate at least len(imgs)

	for _, img := range imgs {
		nfvImgs = append(nfvImgs, cln.makeImage(img)...)
	}

	return nfvImgs
}

// makeNetwork converts a pop Network into a catalogue Network.
func (cln *Client) makeNetwork(net *pop.Network) *catalogue.Network {
	subs := make([]*catalogue.Subnet, len(net.Subnets))

	for i, rsub := range net.Subnets {
		subs[i] = &catalogue.Subnet{
			ExtID:     net.Id,
			CIDR:      rsub.Cidr,
			GatewayIP: rsub.Gateway,
		}
	}

	return &catalogue.Network{
		ExtID:    net.Id,
		Name:     net.Name,
		External: net.External,
		Subnets:  subs,
	}
}

// makeNetwork converts a list of pop Network into a list of catalogue Network.
func (cln *Client) makeNetworks(rnets []*pop.Network) []*catalogue.Network {
	nets := make([]*catalogue.Network, len(rnets))

	for i, rnet := range rnets {
		nets[i] = cln.makeNetwork(rnet)
	}

	return nets
}

// makeServer converts a pop Container into a catalogue Server.
func (cln *Client) makeServer(ctx context.Context, cont *pop.Container) (srv *catalogue.Server, err error) {
	// also fetch the image
	var nfvImage *catalogue.NFVImage
	if cont.ImageId != "" {
		nfvImage, err = cln.Image(ctx, IDFilter(cont.ImageId))
		if err != nil {
			return nil, err
		}
	}

	var deploymentFlavour *catalogue.DeploymentFlavour
	if cont.FlavourId != "" {
		deploymentFlavour, err = cln.Flavour(ctx, IDFilter(cont.FlavourId))
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
		for netID, ep := range cont.Endpoints {
			if ep.Ipv4 != nil && ep.Ipv4.Address != "" {
				name := netID
				if ep.NetName != "" {
					name = ep.NetName
				}

				// no ipv6 for now
				ipMap[name] = []string{ep.Ipv4.Address}
			}
		}
	}

	return &catalogue.Server{
		ExtID:          cont.Id,
		Name:           name,
		Status:         cont.Status.String(),
		ExtendedStatus: cont.ExtendedStatus,
		Image:          nfvImage,
		Flavour:        deploymentFlavour,
		IPs:            ipMap,
		FloatingIPs:    map[string]string{},
		Created:        catalogue.UnixDate(cont.Created),
	}, nil
}

// makeServers converts a list of pop Containers into a list of catalogue Server.
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
