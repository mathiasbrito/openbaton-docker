package client

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

//go:generate protoc -I ../proto ../proto/pop.proto --go_out=plugins=grpc:../proto

// Client is a client instance for a Pop, that automatically converts
// Pop-protocol values into OpenBaton catalogue types.
// Clients use cached connections, and they are identified by their Credentials.
type Client struct {
	Credentials creds.Credentials
}

// New returns a Client initialized with credentials extracted from a given VIMInstance.
func New(inst *catalogue.VIMInstance) *Client {
	c := creds.FromVIM(inst)

	return &Client{Credentials: c}
}

// Info retrieves informations from the current PoP.
func (cln *Client) Info(ctx context.Context) (infos *pop.Infos, err error) {
	err = cln.doRetry(func(stub pop.PopClient) (err error) {
		infos, err = stub.Info(ctx, &empty.Empty{})
		return
	})

	return
}

// FetchMetadata fetches the metadata for a given server.
// This function is generally not needed by a normal user of this library.
func (cln *Client) FetchMetadata(ctx context.Context, f Filter) (md map[string]string, err error) {
	err = cln.doRetry(func(stub pop.PopClient) error {
		conts, err := stub.Containers(ctx, filter(f))
		if err != nil {
			return err
		}

		if conts == nil || conts.List == nil || len(conts.List) != 1 {
			return fmt.Errorf("invalid argument returned from server: %s", conts.String())
		}

		if protoMd := conts.List[0].Md; protoMd != nil && protoMd.Entries != nil {
			md = protoMd.Entries
		} else {
			md = map[string]string{}
		}
		
		return nil
	})

	return
}

// sessionOp is the type of the callback of doRetry.
type sessionOp func(pop.PopClient) error

// doRetry is an helper method that executes an RPC call, retrying it in case
// the session becomes invalid.
func (cln *Client) doRetry(op sessionOp) error {
	for {
		// In case there's no currently cached stub,
		// or if the stub is invalid, a new one will be created by logging into
		// the service again.
		stub, err := cln.stub()
		if err != nil {
			return err
		}

		// if the error is nil or not from an invalid token, do this again
		if err := op(stub); err != errInvalidSession {
			return err
		}
	}
}

// stub() is a helper method to make retrieving a stub from the cache quicker.
func (cln *Client) stub() (pop.PopClient, error) {
	sess, err := cache.get(cln.Credentials)
	if err != nil {
		return nil, err
	}

	return sess.stub(), nil
}

// Filter represents a filter type to be applied during a server query.
type Filter interface{
	isFilterType() // dummy method to force users to use the filters below
}

// IDFilter contains the ID that should be matched by an operation.
type IDFilter string

func (IDFilter) isFilterType() {}

// NameFilter contains the name that will be matched during an operation.
type NameFilter string

func (NameFilter) isFilterType() {}

func filter(cf Filter) *pop.Filter {
	filter := &pop.Filter{}

	switch f := cf.(type) {
	case IDFilter: 
		filter.Options = &pop.Filter_Id{
			Id: string(f),
		}

	case NameFilter:
		filter.Options = &pop.Filter_Name{
			Name: string(f),
		}

	default:
		panic("unknown filter type")
	}

	return filter
}