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
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
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

	// the container should pass through his events only once.
	mux sync.Mutex
}

func (pcont *svcCont) Md() metadata {
	return pcont.ToPop().Md.Entries
}

func (pcont *svcCont) ToPop() *pop.Container {
	return pcont.Container
}

// concrete service
type service struct {
	*log.Logger
	sessionManager
	users Users
	name  string
	cln   *client.Client
	conts map[string]*svcCont

	nets       map[string]*svcNet
	defaultNet *svcNet // also in nets

	// contNames is a map of name -> id for conts;
	// this allows fast matching of the id from the name
	contNames map[string]string

	// same for nets
	netNames map[string]string

	contsMux sync.RWMutex
	quitChan chan struct{}
}

func newService(cfg Config, l *log.Logger) (*service, error) {
	tag := util.FuncName()

	l.WithFields(log.Fields{
		"tag": tag,
	}).Debug("creating route service")

	cln, err := dialDocker(cfg)
	if err != nil {
		return nil, err
	}

	svc := &service{
		Logger: l,

		name: cfg.PopName,
		cln:  cln,
		sessionManager: sessionManager{
			tk: make(map[string]struct{}),
		},

		users:     cfg.Users,
		conts:     make(map[string]*svcCont),
		nets:      make(map[string]*svcNet),
		contNames: make(map[string]string),
		netNames:  make(map[string]string),
		quitChan:  make(chan struct{}),
	}

	if err := svc.checkDocker(); err != nil {
		return nil, fmt.Errorf("docker connection is broken: %v", err)
	}

	l.WithFields(log.Fields{
		"tag": tag,
	}).Debug("creating private network if not present...")

	svc.defaultNet, err = svc.initPrivateNetwork(defaultNetName)
	if err != nil {
		return nil, err
	}

	l.WithFields(log.Fields{
		"tag":        tag,
		"net-name":   defaultNetName,
		"net-subnet": svc.defaultNet.net4,
	}).Debug("obtained default private network")

	// spawn the monitoring loop
	go svc.refreshLoop()

	return svc, nil
}

func (svc *service) Info(context.Context, *empty.Empty) (*pop.Infos, error) {
	tag := util.FuncName()
	op := "Info"

	svc.WithFields(log.Fields{
		"tag": tag,
		"op":  op,
	}).Debug("infos requested")

	return &pop.Infos{
		Name:      svc.name,
		Type:      "docker",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (svc *service) checkDocker() (err error) {
	tag := util.FuncName()

	svc.WithFields(log.Fields{
		"tag": tag,
	}).Debug("checking Docker daemon")

	_, err = svc.cln.Ping(context.Background())
	return
}

func (svc *service) close() error {
	tag := util.FuncName()

	svc.WithFields(log.Fields{
		"tag": tag,
	}).Debug("stopping route service")

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
