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
	net.Listener
}

func New(confFile string) (*Server, error) {
	cfg, err := ReadConfigFile(confFile)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config: cfg,
	}, nil
}

// Serve spawns the service.
func (s *Server) Serve() error {
	log.WithField("cfg", s.Config).Info("starting docker-popd")

	lis, err := net.Listen(s.Config.Proto, s.Config.Netaddr)
	if err != nil {
		return err
	}

	s.Listener = lis

	svc, err := newService(s.Config)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(svc.streamInterceptor),
		grpc.UnaryInterceptor(svc.unaryInterceptor),
	)

	pop.RegisterPopServer(srv, svc)

	if err := srv.Serve(s.Listener); err != nil {
		return err
	}

	return nil
}
