package server

import (
	"net"

	"google.golang.org/grpc"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

// Server represents the PoP service.
type Server struct {
	Config
	net.Listener
}

// New initialises a new Server from a configuration file.
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
	//log.WithField("cfg", s.Config).Info("starting docker-popd")

	proto := s.Config.Proto
	if proto == "" {
		proto = pop.DefaultListenProtocol
	}
	
	laddr := s.Config.Netaddr
	if laddr == "" {
		laddr = pop.DefaultListenAddress
	}

	lis, err := net.Listen(proto, laddr)
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
