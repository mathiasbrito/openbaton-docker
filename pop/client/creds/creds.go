// Package creds implements a Credentials structure to represent and store Pop credentials.
// It also provides functions to extract credentials from a VIMInstance.
package creds

import (
	"net/url"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/catalogue"
)

// Credentials to connect and authenticate with a Pop server.
type Credentials struct {
	Host     string
	Username string
	Password string
}

// FromVIM extracts pop credentials stored in a VIMInstance.
func FromVIM(vimInstance *catalogue.VIMInstance) Credentials {
	host := vimInstance.AuthURL

	// try to parse the AuthURL as an URL (it should be).
	// In case this fails, our last resort will be to use it as an host, hoping for the best;
	// if it is broken, it will fail later during the first request.
	hostURL, err := url.Parse(vimInstance.AuthURL)
	if err == nil {
		host = hostURL.Host
	}

	return Credentials{
		Host:     host,
		Username: vimInstance.Username,
		Password: vimInstance.Password,
	}
}

// ToPop converts the Credentials to a pop.Credentials instance.
func (c *Credentials) ToPop() *pop.Credentials {
	return &pop.Credentials{
		Username: c.Username,
		Password: c.Password,
	}
}
