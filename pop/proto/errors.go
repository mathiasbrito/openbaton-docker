package proto

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	AuthErr         = grpc.Errorf(codes.InvalidArgument, "invalid credentials")
	InternalErr     = grpc.Errorf(codes.Internal, "server fault")
	InvalidArgErr   = grpc.Errorf(codes.InvalidArgument, "invalid arguments")
	InvalidTokenErr = grpc.Errorf(codes.PermissionDenied, "invalid token")
	NotLoggedErr    = grpc.Errorf(codes.Unauthenticated, "not authenticated")
)
