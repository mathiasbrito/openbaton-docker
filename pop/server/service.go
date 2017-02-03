package server

import (
	"errors"
	"fmt"
	"sync"
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
)

type metadata map[string]string

// Merge merges metadata maps together.
// New values will overwrite old ones, and empty values
// will delete the key from the map.
func (md metadata) Merge(newMD metadata) {
	for key, val := range newMD {
		if val != "" {
			md[key] = val
		} else {
			delete(md, key)
		}
	}
}

func (md metadata) Strings() []string {
	ret := make([]string, 0, len(md))

	for key, val := range md {
		ret = append(ret, fmt.Sprintf("%s=%s", key, val))
	}

	return ret
}

// svcCont represent a link between a Pop Container
// and a Docker container.
type svcCont struct {
	*pop.Container
	DockerID string

	metadata

	// the container should pass through his events only once.
	mux sync.Mutex
}

// concrete service
type service struct {
	sessionManager
	users    Users
	name     string
	cln      *client.Client
	conts    map[string]*svcCont
	contsMux sync.RWMutex
	quitChan chan struct{}
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
		users:    cfg.Users,
		conts:    make(map[string]*svcCont),
		quitChan: make(chan struct{}),
	}

	if err := srv.checkDocker(); err != nil {
		return nil, fmt.Errorf("docker connection is broken: %v", err)
	}

	// spawn the monitoring loop
	go srv.refreshLoop()

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

func (svc *service) close() error {
	svc.quitChan <- struct{}{}

	select {
	case <-svc.quitChan:
		return nil

	case <-time.After(5 * time.Second):
		return errors.New("timed out while closing the Docker monitor routine")
	}
}

func dialDocker(cfg Config) (*client.Client, error) {
	host := cfg.DockerdHost
	if host == "" {
		host = client.DefaultDockerHost
	}

	return client.NewClient(host, client.DefaultVersion, nil, nil)
}
