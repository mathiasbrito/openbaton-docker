package server

import (
	"crypto/rand"
	"encoding/base64"
	"sync"

	"golang.org/x/net/context"

	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/openbaton-docker/pop"
)

type sessionManager struct {
	l  sync.RWMutex
	tk map[string]struct{}
}

func (sm *sessionManager) CheckToken(tok string) bool {
	sm.l.RLock()
	defer sm.l.RUnlock()

	_, ok := sm.tk[tok]
	return ok
}

func (sm *sessionManager) DeleteToken(tok string) {
	delete(sm.tk, tok)
}

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
	if creds == nil {
		return nil, pop.InvalidArgErr
	}

	if user, found := svc.users[creds.Username]; found {
		if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(creds.Password)) == nil {
			tok, err := svc.NewToken()
			if err != nil {
				return nil, pop.InternalErr
			}

			return &pop.Token{Value: tok}, nil
		}
	}

	return nil, pop.AuthErr
}

func (svc *service) Logout(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	// getToken() will always return a valid token (it has been checked in unaryInterceptor()).

	for _, token := range getTokens(ctx) {
		svc.DeleteToken(token)
	}

	return &empty.Empty{}, nil
}

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

func (svc *service) streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := svc.authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

func (svc *service) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Let the Login method AND ONLY IT pass through without a valid token (for obvious reasons)
	if info.FullMethod != loginMethod {
		if err := svc.authorize(ctx); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func getTokens(ctx context.Context) []string {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return []string{}
	}

	if len(md["token"]) == 0 {
		return []string{}
	}

	return md["token"]
}
