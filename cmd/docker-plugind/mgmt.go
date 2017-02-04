package main

import (
	"context"
	"errors"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/openbaton-docker/mgmt"
	"github.com/mcilloni/openbaton-docker/pop/client"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
)

var (
	// ErrMgmtUnavailable indicates that the plugin has been unable to spawn management for the given VIM instance.
	ErrMgmtUnavailable = errors.New("management is unavailable for this instance")
)

func (d *driver) SetupManagement(vimInstance *catalogue.VIMInstance) (bool, error) {
	if d.Accessor == nil {
		return false, ErrMgmtUnavailable
	}

	id := vimInstance.ID

	// If there is no VIMInstance ID, then finding the VIM is impossible.
	if id == "" {
		return false, nil
	}

	// if the manager for the given VIMInstance is already on,
	// then don't do anything, we are already set up
	if _, on := d.managers[id]; on {
		return false, nil
	}

	c := creds.FromVIM(vimInstance)
	d.managers[id] = mgmt.NewManager(id, newHandler(c), d.Accessor, d.Logger)

	return true, nil
}

type handler struct {
	cln client.Client
}

// newHandler created an handler with its Credentials stored in a client.Client instance.
func newHandler(c creds.Credentials) handler {
	return handler{client.Client{Credentials: c}}
}

func (h handler) AddMetadata(name string, entries map[string]string) error {
	return h.cln.AddMetadata(context.Background(), client.NameFilter(name), entries)
}

// Check checks on the pop if the Server with the given name exists.
func (h handler) Check(name string) (*catalogue.Server, error) {
	return h.cln.Server(context.Background(), client.NameFilter(name))
}

// Start starts the given server.
func (h handler) Start(name string) error {
	_, err := h.cln.Start(context.Background(), client.NameFilter(name))

	return err
}
