package server

import (
    "context"
    "errors"
    "fmt"
    "net"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
    pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/docker/docker/api/types/network"
)

// this file implements a very very basic IP allocation system, that assigns 
// to the private network incremental IPs. It does not support IP reassignment after
// use.

const (
    privateNetName = "popd-private"
)

var (
    ErrIPExhausted = errors.New("no more IP available for this network")
    ErrSubnetsExhausted = errors.New("no more subnets available for this host")
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
                IP: sub,
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

func (svc *service) getPrivateEndpoint() (*pop.Endpoint, error) {
    ip4, _, err := svc.privNet.GetV4()
    if err != nil {
        return nil, err
    }

    return &pop.Endpoint{
        NetId: svc.privNet.ID,
        Ipv4: &pop.Ip{Address: ip4.String()},
        Ipv6: &pop.Ip{},
    }, nil
}

func (svc *service) initPrivateNetwork() error {
    net, err := svc.cln.NetworkInspect(context.Background(), privateNetName)

    // If the network doesn't exist, create a new one.
    if client.IsErrNetworkNotFound(err) {
        net, err = svc.newPrivateNetwork()
    } 
    
    // some other error happened
    if err != nil {
        return err
    }

    svc.privNet, err = newSvcNet(net)
    if err != nil {
        return err
    }

    return nil
}

func (svc *service) newPrivateNetwork() (nr types.NetworkResource, opErr error) {
    // create a new private network

    ipNet, err := svc.detectNewSubnet4()
    if err != nil {
        opErr = err
        return 
    }

    opts := types.NetworkCreate{
        Attachable: true,
        Internal: true,
        IPAM: &network.IPAM{
            Config: []network.IPAMConfig{
                {
                    Subnet: ipNet.String(),
                },
            },
        },
    }

    resp, err := svc.cln.NetworkCreate(context.Background(), privateNetName, opts)
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

type svcNet struct {
    ID string
    current4 net.IP
    gateway4 net.IP
    net4     *net.IPNet   
    exhausted bool
}

func newSvcNet(dnr types.NetworkResource) (pnet svcNet, opErr error) {
    if len(dnr.IPAM.Config) < 1 {
        opErr = fmt.Errorf("malformed network %s", dnr.Name)
        return
    }

    pnet.ID = dnr.ID

    // I just need a subnet; the first one will be ok
    subnet4 := dnr.IPAM.Config[0].Subnet
    if subnet4 == "" {
        opErr = fmt.Errorf("network %s has no IPv4 subnet", dnr.Name)
        return
    }

    // load the subnet, and use it as the base IP of the pnet
    pnet.current4, pnet.net4, opErr = net.ParseCIDR(subnet4)
    if opErr != nil {
        return
    }

    // add the gateway if present
    gateway4 := dnr.IPAM.Config[0].Gateway
    if gateway4 != "" {
        pnet.gateway4 = net.ParseIP(gateway4)
    }

    return
}

func (pnet *svcNet) GetV4() (net.IP, net.IPMask, error) {
    if err := pnet.nextV4(); err != nil {
        return nil, nil, err
    }

    if pnet.current4.Equal(pnet.gateway4) {
        // skip the gateway and get the next one
        return pnet.GetV4() 
    }

    return pnet.current4, pnet.net4.Mask, nil
}

func (pnet *svcNet) nextV4() error {
    if pnet.exhausted {
        return ErrIPExhausted
    }

    // an IP Address is big endian
    for i := len(pnet.current4) - 1; i >= 0; i-- {
        pnet.current4[i]++

        // if there is no overflow, stop
        if pnet.current4[i] != 0 {
            break
        }
    }

    // notify the caller that the network is exhausted
    if !pnet.net4.Contains(pnet.current4) {
        pnet.exhausted = true

        return ErrIPExhausted
    }

    return nil
}

func inIPSlice(haystack []net.IP, needle net.IP) bool {
    for _, hip := range haystack {
        if needle.Equal(hip) {
            return true
        }
    }

    return false
}