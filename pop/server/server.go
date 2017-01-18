package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/mcilloni/openbaton-docker/pop"
	log "github.com/sirupsen/logrus"
)

// Server represents the PoP service.
type Server struct {
	Config
}

func New(confFile string) (Server, error) {
	cfg, err := ReadConfigFile(confFile)
	if err != nil {
		return Server{}, err
	}

	return Server{cfg}, nil
}

// Serve spawns the service.
func (s Server) Serve() error {
	log.WithField("cfg", s.Config).Info("starting docker-popd")

	lis, err := net.Listen(s.Config.Proto, s.Config.Netaddr)
	if err != nil {
		return err
	}

	svc, err := newService(s.Config)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(svc.streamInterceptor),
		grpc.UnaryInterceptor(svc.unaryInterceptor),
	)

	pop.RegisterPoPServer(srv, svc)

	if err := srv.Serve(lis); err != nil {
		return err
	}

	return nil
}
