package client

import (
	"errors"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/openbaton-docker/pop/client/creds"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

var (
	// errInvalidSession signals that the session is not valid anymore.
	errInvalidSession = errors.New("the client is invalid; retry again")
)

// connection represents a session with the server.
// connection instances are cached and discarded in case they become invalid.
type session struct {
	conn *grpc.ClientConn

	tok string

	invalid bool
}

// newSession initialises a session, authenticating into the service
// and getting a token.
func newSession(creds creds.Credentials) (*session, error) {
	sess := new(session)

	target := creds.Host
	if target == "" {
		target = pop.DefaultAddress
	}

	// WithInsecure allows for non-TLS connections.
	gconn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(sess.interceptor),
	)

	if err != nil {
		return nil, err
	}

	sess.conn = gconn

	s := sess.stub()

	// create a new session by logging into the service
	tk, err := s.Login(context.Background(), creds.ToPop())
	if err != nil {
		return nil, err
	}

	// store the token
	sess.tok = tk.Value

	return sess, nil
}

// ctx returns a Context in which the token has been set as metadata.
func (sess *session) ctx(ctx context.Context) context.Context {
	return metadata.NewContext(ctx, metadata.Pairs(pop.TokenKey, sess.tok))
}

// interceptor intercepts each call, injects the token, executes the call and then checks if the token is valid.
// In case  it is invalid, it marks the current session as invalid.
func (sess *session) interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if sess.invalid {
		return errInvalidSession
	}

	// if we are not logging in, inject the token metadata in the context
	if method != "/pop.Pop/Login" {
		ctx = sess.ctx(ctx)
	}

	err := invoker(ctx, method, req, reply, cc, opts...)

	// If the error we got is permission denied, the token is not valid, so the connection structure
	// must be dropped.
	if grpc.Code(err) == codes.PermissionDenied {
		sess.invalid = true

		return errInvalidSession
	}

	return err
}

// logout logs the session out of the service, invalidating it.
func (sess *session) logout() error {
	stub := sess.stub()

	_, err := stub.Logout(context.Background(), &empty.Empty{})

	// logging out invalids the session
	sess.invalid = true

	if err != nil && err != errInvalidSession {
		return err
	}

	return nil
}

// stub uses the underlying connection to create a new stub.
func (sess *session) stub() pop.PopClient {
	return pop.NewPopClient(sess.conn)
}
