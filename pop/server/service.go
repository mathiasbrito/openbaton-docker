package server

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

//go:generate protoc -I ../proto ../proto/pop.proto --go_out=plugins=grpc:../proto

const (
	// TokenBytes specifies how long a token is.
	TokenBytes = 32

	// loginMethod is the signature of the login method. Check this string 
	// carefully.
	loginMethod = "/pop.Pop/Login"
)

// concrete service 
type service struct {
	sessionManager
	users Users
	name  string
	cln   *client.Client
}

func newService(cfg Config) (*service, error) {
	cln, err := dialDocker(cfg)
	if err != nil {
		return nil, err
	}

	srv := &service{
		name: cfg.PopName,
		cln:  cln,
		sessionManager: sessionManager{
			tk: make(map[string]struct{}),
		},
		users: cfg.Users,
	}

	if err := srv.checkDocker(); err != nil {
		return nil, fmt.Errorf("docker connection is broken: %v", err)
	}

	return srv, nil
}

func (svc *service) Info(context.Context, *empty.Empty) (*pop.Infos, error) {
	return &pop.Infos{
		Name:      svc.name,
		Type:      "docker",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (svc *service) checkDocker() (err error) {
	_, err = svc.cln.Ping(context.Background())
	return
}

func dialDocker(cfg Config) (*client.Client, error) {
	host := cfg.DockerdHost
	if host == "" {
		host = client.DefaultDockerHost
	}

	return client.NewClient(host, client.DefaultVersion, nil, nil)
}
