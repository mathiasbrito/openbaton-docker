package server

import (
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/client"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mcilloni/openbaton-docker/pop"
)

//go:generate protoc -I ../proto ../proto/pop.proto --go_out=plugins=grpc:..

const (
	// TokenBytes specifies how long a token is.
	TokenBytes = 32

	loginMethod = "/pop.PoP/Login"
)

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

	return &service{
		name: cfg.PopName,
		cln:  cln,
		sessionManager: sessionManager{
			tk: make(map[string]struct{}),
		},
		users: cfg.Users,
	}, nil
}

func (svc *service) Info(context.Context, *empty.Empty) (*pop.Infos, error) {
	return &pop.Infos{
		Name:      svc.name,
		Type:      "docker",
		Timestamp: time.Now().Unix(),
	}, nil
}

func dialDocker(cfg Config) (*client.Client, error) {
	host := cfg.DockerdHost
	if host == "" {
		host = client.DefaultDockerHost
	}

	return client.NewClient(host, client.DefaultVersion, nil, nil)
}
