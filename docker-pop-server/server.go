package server

import (
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pop "github.com/mcilloni/openbaton-docker/pop/proto"
	"github.com/openbaton/go-openbaton/util"
	log "github.com/sirupsen/logrus"
)

// Server represents the PoP service.
type Server struct {
	*log.Logger
	Config   Config
	GRPCServer *grpc.Server

	svc *service
}

// New initialises a new Server from viper.
func New() (*Server, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	l := newLogger(cfg)

	return &Server{
		Logger: l,
		Config: cfg,
	}, nil
}

// Close shuts down the Server.
func (s *Server) Close() error {
	tag := util.FuncName()

	s.WithFields(log.Fields{
		"tag":      tag,
		"pop-name": s.Config.PopName,
	}).Info("stopping server")

	s.GRPCServer.Stop()
	
	return s.svc.close()
}

// Serve spawns the service.
func (s *Server) Serve() error {
	tag := util.FuncName()

	s.WithFields(log.Fields{
		"tag":      tag,
		"pop-name": s.Config.PopName,
	}).Info("starting server")

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

	s.svc, err = newService(s.Config, s.Logger)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(s.svc.streamInterceptor),
		grpc.UnaryInterceptor(s.svc.unaryInterceptor),
	}

	if s.Config.TLSCertPath != "" && s.Config.TLSKeyPath != "" {
		creds, err := credentials.NewServerTLSFromFile(s.Config.TLSCertPath, s.Config.TLSKeyPath)
		if err != nil {
			return err
		}

		opts = append(opts, grpc.Creds(creds))
	}

	s.GRPCServer = grpc.NewServer(opts...)


	pop.RegisterPopServer(s.GRPCServer, s.svc)

	s.WithFields(log.Fields{
		"tag":      tag,
		"pop-name": s.Config.PopName,
	}).Info("launching gRPC server")

	err = s.GRPCServer.Serve(lis)

	if got, want := grpc.ErrorDesc(err), "use of closed network connection"; got != "" && !strings.Contains(got, want) {
		return err
	}

	return nil
}
