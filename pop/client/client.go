package client

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/go-openbaton/catalogue"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
)

//go:generate protoc -I ../proto ../proto/pop.proto --go_out=plugins=grpc:../proto

// Client is a client instance for a Pop, that automatically converts
// Pop-protocol values into OpenBaton catalogue types.
// Clients use cached connections, and they are identified by their Credentials.
type Client struct {
	Credentials creds.Credentials
}

// New returns a Client for given instance, initializing it with
// credentials extracted from the given VIMInstance.
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

func (cln *Client) stub() (pop.PopClient, error) {
	sess, err := cache.get(cln.Credentials)
	if err != nil {
		return nil, err
	}

	return sess.stub(), nil
}
