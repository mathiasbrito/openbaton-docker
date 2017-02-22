package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

// this file implements a very very basic IP allocation system, that takes care of
// which IPs were allocated on the network.

const (
	defaultNetName = "private"
)

var (
	ErrAddrTaken        = grpc.Errorf(codes.AlreadyExists, "ip is not available")
	ErrIPExhausted      = grpc.Errorf(codes.ResourceExhausted, "no more IP available for this network")
	ErrSubnetsExhausted = grpc.Errorf(codes.ResourceExhausted, "no more subnets available for this host")
)

// detectNewSubnet scans every subnet known by docker, finds their subnets, and then finds a
// suitable, free subnet to use.
func (svc *service) detectNewSubnet4() (*net.IPNet, error) {
	taken, err := svc.fetchDockerSubnets4()
	if err != nil {
		return nil, err
	}

	// A /12 is reserved for private use, from 172.16.0.0 to 172.31.255.255 (obv w.o. .0 and .255).
	// This means that there are 2**4 -> 16 available subnets in that range;
	// Docker often uses these to make its NAT networks, so I guess it's fine to use one of those too
	// to get a /16.

	// parentheses to make the syntax not ambiguous
	for sub := (net.IP{172, 16, 0, 0}); sub[1] < 32; sub[1]++ {
		// this subnet is unused
		if !inIPSlice(taken, sub) {
			return &net.IPNet{
				IP:   sub,
				Mask: net.IPv4Mask(255, 255, 0, 0), // /16
			}, nil
		}
	}

	return nil, ErrSubnetsExhausted
}

// fetches only subnets between 172.16.0.0 and 172.31.0.0
func (svc *service) fetchDockerSubnets4() ([]net.IP, error) {
	dnets, err := svc.cln.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	ret := []net.IP{}

	for _, dnet := range dnets {
		for _, cfg := range dnet.IPAM.Config {
			_, ipNet, err := net.ParseCIDR(cfg.Subnet)
			if err != nil {
				continue
			}

			// if it is not a /16 network, ignore it
			if ones, _ := ipNet.Mask.Size(); ones != 16 {
				continue
			}

			ip := ipNet.IP.To4()

			// ignore IPv6 here
			if ip == nil {
				continue
			}

			// check that the subnet IP is in the right range.
			if ip[0] != 172 || ip[1] < 16 || ip[1] > 31 {
				continue
			}

			ret = append(ret, ip)
		}
	}

	return ret, nil
}

func (svc *service) getDefaultEndpoint() (*pop.Endpoint, error) {
	ip4, _, err := svc.defaultNet.GetV4()
	if err != nil {
		return nil, err
	}

	return &pop.Endpoint{
		NetId: svc.defaultNet.ID,
		Ipv4:  &pop.Ip{Address: ip4.String()},
		Ipv6:  &pop.Ip{},
	}, nil
}

func (svc *service) initPrivateNetwork(name string) (*svcNet, error) {
	net, err := svc.cln.NetworkInspect(context.Background(), name)

	// If the network doesn't exist, create a new one.
	if client.IsErrNetworkNotFound(err) {
		net, err = svc.newRandomNetwork(name)
	}

	// some other error happened
	if err != nil {
		return nil, err
	}

	newNet, err := newSvcNet(name, net)
	if err != nil {
		return nil, err
	}

	svc.nets[newNet.ID] = newNet
	svc.netNames[name] = newNet.ID

	return newNet, nil
}

func (svc *service) newRandomNetwork(name string) (nr types.NetworkResource, opErr error) {
	// create a new network using a free /16

	ipNet4, err := svc.detectNewSubnet4()
	if err != nil {
		opErr = err
		return
	}

	gateway4 := ipNet4.IP
	gateway4[len(gateway4)-1] = 1 // give  the network a standard gateway on address .1

	opts := types.NetworkCreate{
		Attachable: true,
		Internal:   true,
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet:  ipNet4.String(),
					Gateway: gateway4.String(),
				},
			},
		},
	}

	resp, err := svc.cln.NetworkCreate(context.Background(), name, opts)
	if err != nil {
		opErr = err
		return
	}

	newNet, err := svc.cln.NetworkInspect(context.Background(), resp.ID)
	if err != nil {
		opErr = err
		return
	}

	return newNet, nil
}

func (svc *service) releaseContIPs(pcont *svcCont) error {
	tag := util.FuncName()

	svc.WithFields(log.Fields{
		"tag":         tag,
		"pcont-names": pcont.Names,
	}).Debug("releasing IPs for container")

	if pcont.Endpoints == nil {
		return nil
	}

	for nid, ep := range pcont.Endpoints {
		if ep.Ipv4 == nil {
			return nil
		}

		ip := net.ParseIP(ep.Ipv4.Address)
		if ip == nil {
			return nil
		}

		pnet := svc.nets[nid]
		if pnet != nil {
			pnet.ReturnV4(ip)

			svc.WithFields(log.Fields{
				"tag":               tag,
				"pcont-names":       pcont.Names,
				"pcont-ip":          ip,
				"net-name":          pnet.Name,
				"net-allocated-ips": len(pnet.taken4),
			}).Debug("private IP released")
		}
	}

	return nil
}

// svcNet wraps and represents a network handled by Popd.
// While Pop represents networks as capable of having multiple subnetworks, this implementation
// only supports a single subnet.
type svcNet struct {
	ID   string
	Name string

	current4 net.IP
	external bool
	gateway4 net.IP
	net4     *net.IPNet
	taken4   map[[4]byte]struct{} // struct{} has 0 size - this is actually a set
	mux      sync.Mutex
}

func newSvcNet(name string, dnr types.NetworkResource) (*svcNet, error) {
	if len(dnr.IPAM.Config) < 1 {
		return nil, fmt.Errorf("malformed network %s", dnr.Name)
	}

	// I just need a subnet; the first one will be ok
	subnet4 := dnr.IPAM.Config[0].Subnet
	if subnet4 == "" {
		return nil, fmt.Errorf("network %s has no IPv4 subnet", dnr.Name)
	}

	// load the subnet, and use it as the base IP of the pnet
	current4, net4, err := net.ParseCIDR(subnet4)
	if err != nil {
		return nil, err
	}

	// shrink the IP to 4 bytes
	current4 = current4.To4()

	taken4, err := parseTaken(dnr.Containers)
	if err != nil {
		return nil, err
	}

	// add the gateway if present
	gateway4 := net.ParseIP(dnr.IPAM.Config[0].Gateway)

	// shrink the IP to 4 bytes
	gateway4 = gateway4.To4()

	return &svcNet{
		ID:   dnr.ID,
		Name: name,

		current4: current4,
		external: !dnr.Internal,
		gateway4: gateway4,
		net4:     net4,
		taken4:   taken4,
	}, nil
}

// GetV4 returns an IPv4 on the svcNet.
func (pnet *svcNet) GetV4() (net.IP, net.IPMask, error) {
	pnet.mux.Lock()
	defer pnet.mux.Unlock()

	// check if there are IPv4s left
	if !pnet.hasNext4() {
		return nil, nil, ErrIPExhausted
	}

	// there is at least a valid ip - this loop should always break
	for {
		// get the next IPv4
		pnet.nextV4()

		if pnet.currentIsValid() {
			b4, err := ip4ToArr(pnet.current4)
			if err != nil {
				return nil, nil, err
			}

			_, found := pnet.taken4[b4]
			if !found {
				// assign the IPv4
				pnet.taken4[b4] = struct{}{}
				return pnet.current4, pnet.net4.Mask, nil
			}
		}
	}
}

func (pnet *svcNet) ReserveV4(ip net.IP) error {
	pnet.mux.Lock()
	defer pnet.mux.Unlock()

	if !pnet.net4.Contains(ip) {
		return grpc.Errorf(codes.InvalidArgument, "%v is not contained in %v", ip, pnet.net4)
	}

	ip4, err := ip4ToArr(ip)
	if err != nil {
		return err
	}

	if _, found := pnet.taken4[ip4]; found {
		return ErrAddrTaken
	}

	pnet.taken4[ip4] = struct{}{}

	return nil
}

// ReturnV4 returns an IPv4 to the network.
// If the address is invalid, this operation returns false.
func (pnet *svcNet) ReturnV4(ip net.IP) bool {
	pnet.mux.Lock()
	defer pnet.mux.Unlock()

	ip4, err := ip4ToArr(ip)
	if err != nil {
		return false
	}

	_, found := pnet.taken4[ip4]
	delete(pnet.taken4, ip4)

	return found
}

func (pnet *svcNet) ToPop() *pop.Network {
	return &pop.Network{
		Id:       pnet.ID,
		External: pnet.external,
		Name:     pnet.Name,
		Subnets: []*pop.Subnet{
			{
				Cidr:    pnet.net4.String(),
				Gateway: pnet.gateway4.String(),
			},
		},
	}
}

func (pnet *svcNet) ToPopSubnet() *pop.Subnet {
	return &pop.Subnet{
		Cidr:    pnet.net4.String(),
		Gateway: pnet.net4.String(),
	}
}

func (pnet *svcNet) currentIsValid() bool {
	if pnet.current4.Equal(pnet.gateway4) {
		return false
	}

	lb := pnet.current4[len(pnet.current4)-1]
	if lb == 0 || lb == 255 {
		return false
	}

	return true
}

// hasNext4 would never work with an IPv6 (it uses 128 bit addresses)
func (pnet *svcNet) hasNext4() bool {
	return uint64(len(pnet.taken4)) < pnet.totalValidIPv4s()
}

// nextV4 advances current4 to the next IPv4, without
// caring about if it is valid.
func (pnet *svcNet) nextV4() {
	// an IP Address is big endian
	for i := len(pnet.current4) - 1; i >= 0; i-- {
		pnet.current4[i]++

		// if there is no overflow, stop
		if pnet.current4[i] != 0 {
			break
		}
	}

	// Overflow - reset the counter
	if !pnet.net4.Contains(pnet.current4) {
		pnet.current4 = pnet.net4.IP.To4()
		pnet.current4[3]++ // skip the first byte
	}
}

func (pnet *svcNet) totalValidIPv4s() uint64 {
	maskBits, _ := pnet.net4.Mask.Size()
	addrBits := uint64(32) - uint64(maskBits)
	maxNumOfAddr := uint64(1) << addrBits // 2 ** addrBits

	// each /24 is 8 bits, so if i.e. there are 16 address bits,
	// then there can be 2 ** (addrBits - 8) /24s
	nOf24 := uint64(1) << (addrBits - uint64(8))

	// there are 2 invalid ips in each /24 (.0 and .255)
	maxNumOfAddr -= nOf24 * 2

	// the gateway must too be subtracted from the
	// available IPs
	if pnet.gateway4 != nil {
		maxNumOfAddr--
	}

	return maxNumOfAddr
}

func inIPSlice(haystack []net.IP, needle net.IP) bool {
	for _, hip := range haystack {
		if needle.Equal(hip) {
			return true
		}
	}

	return false
}

func ip4ToArr(ip net.IP) (ipArr [4]byte, opErr error) {
	ip4 := ip.To4()
	if ip4 == nil {
		opErr = fmt.Errorf("%v is not an IPv4", ip)
		return
	}

	copy(ipArr[:], ip4)

	return
}

// parseTaken creates a set of taken IPs already found on the network.
func parseTaken(dconts map[string]types.EndpointResource) (map[[4]byte]struct{}, error) {
	taken4 := make(map[[4]byte]struct{})

	for _, dcont := range dconts {
		if dcont.IPv4Address != "" {
			ip, _, err := net.ParseCIDR(dcont.IPv4Address)
			if err != nil {
				continue
			}

			ip4, err := ip4ToArr(ip)
			if err != nil {
				return nil, err
			}

			taken4[ip4] = struct{}{}
		}
	}

	return taken4, nil
}
