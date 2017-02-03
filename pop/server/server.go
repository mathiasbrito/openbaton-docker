package server

import (
	"net"

	"google.golang.org/grpc"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
)

// Server represents the PoP service.
type Server struct {
	Config   Config
	Listener net.Listener

	svc *service
}

// New initialises a new Server from viper.
func New() (*Server, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return &Server{
		Config: cfg,
	}, nil
}

// Close shuts down the Server.
func (s *Server) Close() error {
	err1 := s.svc.close()
	err2 := s.Listener.Close()

	switch {
	case err1 != nil:
		return err1

	case err2 != nil:
		return err2

	default:
		return nil
	}
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

	s.svc, err = newService(s.Config)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(s.svc.streamInterceptor),
		grpc.UnaryInterceptor(s.svc.unaryInterceptor),
	)

	pop.RegisterPopServer(srv, s.svc)

	if err := srv.Serve(s.Listener); err != nil {
		return err
	}

	return nil
}
