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
	return s.Listener.Close()
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
