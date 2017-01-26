package main

import (
	"context"
	"errors"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/util"
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/mgmt"
	log "github.com/sirupsen/logrus"
)

// driver for the Docker plugin.
type driver struct {
	*log.Logger
	Accessor mgmt.AMQPChannelAccessor
	
	managers map[string]mgmt.VIMManager	
}

func newDriver() *driver {
	return &driver{managers: make(map[string]mgmt.VIMManager)}
}

func (d *driver) AddFlavour(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return deploymentFlavour, nil
}

func (d *driver) AddImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageFile []byte) (*catalogue.NFVImage, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return image, nil
}

func (d *driver) AddImageFromURL(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageURL string) (*catalogue.NFVImage, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return image, nil
}

func (d *driver) CopyImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage, imageFile []byte) (*catalogue.NFVImage, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return image, nil
}

func (d *driver) CreateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return network, nil
}

func (d *driver) CreateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return subnet, nil
}

func (d *driver) DeleteFlavour(vimInstance *catalogue.VIMInstance, extID string) (bool, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return true, nil
}

func (d *driver) DeleteImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (bool, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return true, nil
}

func (d *driver) DeleteNetwork(vimInstance *catalogue.VIMInstance, extID string) (bool, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return true, nil
}

func (d *driver) DeleteServerByIDAndWait(vimInstance *catalogue.VIMInstance, id string) error {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
		"server-id": id,
	}).Debug("received request")

	return client.New(vimInstance).Delete(context.Background(), id)
}

func (d *driver) DeleteSubnet(vimInstance *catalogue.VIMInstance, existingSubnetExtID string) (bool, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return true, nil
}

// it's pointless to "wait" for a container creation - it's quick, and it doesn't make sense
// to have different behaviours.
func (d *driver) LaunchInstance(
	vimInstance *catalogue.VIMInstance,
	hostname, image, flavour, keyPair string,
	networks, securityGroups []string,
	userData string) (*catalogue.Server, error) {

	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return d.LaunchInstanceAndWaitWithIPs(vimInstance, hostname, image, flavour, keyPair, networks, securityGroups, userData, nil, nil)
}

func (d *driver) LaunchInstanceAndWait(
	vimInstance *catalogue.VIMInstance,
	hostname, image, flavour, keyPair string,
	networks, securityGroups []string,
	userData string) (*catalogue.Server, error) {

	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return d.LaunchInstanceAndWaitWithIPs(vimInstance, hostname, image, flavour, keyPair, networks, securityGroups, userData, nil, nil)
}

func (d *driver) LaunchInstanceAndWaitWithIPs(
	vimInstance *catalogue.VIMInstance,
	hostname, image, flavour, keyPair string,
	networks, securityGroups []string,
	userData string,
	floatingIps map[string]string,
	keys []*catalogue.Key) (*catalogue.Server, error) {

	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
		"hostname": hostname,
		"image": image,
		"flavour": flavour,
	}).Debug("received request")

	_, err := d.SetupManagement(vimInstance)
	if err != nil {
		d.WithError(err).WithFields(log.Fields{
			"tag": tag,
			"hostname": hostname,
			"image": image,
			"flavour": flavour,
		}).Error("management failed")

		return nil, errors.New("VIM Management refused to start. See the plugin log")
	}

	ips := map[string]string{}
	for _, n := range networks {
		ips[n] = "" // gets a random IP
	}

	return client.New(vimInstance).Create(context.Background(), hostname, image, flavour, ips)
}

func (d *driver) ListFlavours(vimInstance *catalogue.VIMInstance) ([]*catalogue.DeploymentFlavour, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return client.New(vimInstance).Flavours(context.Background())
}

func (d *driver) ListImages(vimInstance *catalogue.VIMInstance) ([]*catalogue.NFVImage, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return client.New(vimInstance).Images(context.Background())
}

func (d *driver) ListNetworks(vimInstance *catalogue.VIMInstance) ([]*catalogue.Network, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return client.New(vimInstance).Networks(context.Background())
}

func (d *driver) ListServer(vimInstance *catalogue.VIMInstance) ([]*catalogue.Server, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return client.New(vimInstance).Servers(context.Background())
}

func (d *driver) NetworkByID(vimInstance *catalogue.VIMInstance, id string) (*catalogue.Network, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return client.New(vimInstance).Network(context.Background(), id)
}

func (d *driver) Quota(vimInstance *catalogue.VIMInstance) (*catalogue.Quota, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return newQuota(), nil
}

func (d *driver) SubnetsExtIDs(vimInstance *catalogue.VIMInstance, networkExtID string) ([]string, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return []string{networkExtID}, nil
}

func (d *driver) Type(vimInstance *catalogue.VIMInstance) (string, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return "docker-pop", nil
}

func (d *driver) UpdateFlavour(vimInstance *catalogue.VIMInstance, deploymentFlavour *catalogue.DeploymentFlavour) (*catalogue.DeploymentFlavour, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return deploymentFlavour, nil
}

func (d *driver) UpdateImage(vimInstance *catalogue.VIMInstance, image *catalogue.NFVImage) (*catalogue.NFVImage, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return image, nil
}

func (d *driver) UpdateNetwork(vimInstance *catalogue.VIMInstance, network *catalogue.Network) (*catalogue.Network, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return network, nil
}

func (d *driver) UpdateSubnet(vimInstance *catalogue.VIMInstance, createdNetwork *catalogue.Network, subnet *catalogue.Subnet) (*catalogue.Subnet, error) {
	tag := util.FuncName()

	d.WithFields(log.Fields{
		"tag": tag,
	}).Debug("received request")

	return subnet, nil
}

func newServer() *catalogue.Server {
	return &catalogue.Server{
		Name:           "server_name",
		ExtID:          "ext_id",
		Created:        catalogue.NewDate(),
		FloatingIPs:    make(map[string]string),
		ExtendedStatus: "ACTIVE",
		Flavour: &catalogue.DeploymentFlavour{
			Disk:       100,
			ExtID:      "ext",
			FlavourKey: "m1.small",
			RAM:        2000,
			VCPUs:      4,
		},
		IPs: make(map[string][]string),
	}
}

func newNetwork(id string) *catalogue.Network {
	return &catalogue.Network{
		Name:  "network_name",
		ID:    id,
		ExtID: "ext_id",
	}
}

func newQuota() *catalogue.Quota {
	return &catalogue.Quota{
		Cores:       99999,
		FloatingIPs: 99999,
		Instances:   99999,
		KeyPairs:    99999,
		RAM:         99999,
		Tenant:      "test-tenant",
	}
}
