package server

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"golang.org/x/net/context"

	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc"
	grpc_md "google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

// sessionManager saves all the generated tokens.
type sessionManager struct {
	l  sync.RWMutex
	tk map[string]struct{}
}

// CheckToken tries to check if the given token is valid.
func (sm *sessionManager) CheckToken(tok string) bool {
	sm.l.RLock()
	defer sm.l.RUnlock()

	_, ok := sm.tk[tok]
	return ok
}

// DeleteToken deletes a token.
func (sm *sessionManager) DeleteToken(tok string) {
	delete(sm.tk, tok)
}

// NewToken reads a TokenBytes long message from a secure RNG, and
// returns it as a base64 token.
func (sm *sessionManager) NewToken() (string, error) {
	b := make([]byte, TokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(b)

	sm.l.Lock()
	defer sm.l.Unlock()

	sm.tk[token] = struct{}{}

	return token, nil
}

// Login logs into the Pop. It should always be the first function called (to setup a token).
// Remember that tokens are transient and not stored, so a new login is needed in case the service dies.
func (svc *service) Login(ctx context.Context, creds *pop.Credentials) (*pop.Token, error) {
	tag := util.FuncName()

	if creds == nil {
		return nil, pop.InvalidArgErr
	}

	svc.WithFields(log.Fields{
		"tag":   tag,
		"creds": *creds,
	}).Debug("received login")

	if user, found := svc.users[creds.Username]; found {
		if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(creds.Password)) == nil {
			tok, err := svc.NewToken()
			if err != nil {
				svc.WithError(err).WithFields(log.Fields{
					"tag": tag,
				}).Error("login error")

				return nil, pop.InternalErr
			}

			return &pop.Token{Value: tok}, nil
		}
	}

	return nil, pop.AuthErr
}

// Logout deletes the current session token from the sessionManager.
func (svc *service) Logout(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	tag := util.FuncName()
	// getTokens() will always return a valid token (it has been checked in unaryInterceptor()).

	for _, token := range getTokens(ctx) {
		svc.WithFields(log.Fields{
			"tag":   tag,
			"token": token,
		}).Debug("logging out")
		svc.DeleteToken(token)
	}

	return &empty.Empty{}, nil
}

// authorize checks if the current context is autheticated (ie, if it contains a valid token).
func (svc *service) authorize(ctx context.Context) error {
	tokens := getTokens(ctx)

	if len(tokens) == 0 {
		return pop.NotLoggedErr
	}

	for _, token := range tokens {
		if svc.CheckToken(token) {
			return nil
		}
	}

	return pop.InvalidTokenErr
}

// streamInterceptor is an interceptor for stream requests.
func (svc *service) streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := svc.authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

// unaryInterceptor intercepts every unary request, and ensures that the caller is authorized before doing anything.
func (svc *service) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Let the Login method AND ONLY IT pass through without a valid token (for obvious reasons)
	if info.FullMethod != pop.LoginMethod {
		if err := svc.authorize(ctx); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// getToken retrieves the tokens from the current context metadata.
func getTokens(ctx context.Context) []string {
	md, ok := grpc_md.FromContext(ctx)
	if !ok {
		return []string{}
	}

	if len(md["token"]) == 0 {
		return []string{}
	}

	return md["token"]
}
