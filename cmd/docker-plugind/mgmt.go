package main

import (
	"context"
	"errors"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	"github.com/mcilloni/openbaton-docker/pop/mgmt"
)

var (
	ErrMgmtUnavailable = errors.New("management is unavailable for this instance")
)

func (d *driver) SetupManagement(vimInstance *catalogue.VIMInstance) (bool, error) {
	if d.accessor == nil {
		return false, ErrMgmtUnavailable
	}

	c := creds.FromVIM(vimInstance)

	// if the manager for the given VIMInstance is already on,
	// then don't do anything, we are already set up
	if _, on := d.managers[c]; on {
		return false, nil
	}

	d.managers[c] = mgmt.NewManager(vimInstance.Name, newHandler(c), d.accessor, d.Logger)

	return true, nil
}

type handler struct {
	cln client.Client
}

func newHandler(c creds.Credentials) handler {
	return handler{client.Client{Credentials: c}}
}

func (h handler) Start(id string) error {
	_, err := h.cln.Start(context.Background(), id)

	return err
}
