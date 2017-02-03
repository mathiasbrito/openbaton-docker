package proto

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	// AuthErr represent an authentication failure (ie, wrong password).
	AuthErr = grpc.Errorf(codes.InvalidArgument, "invalid credentials")

	// InternalErr represents an internal crash of the server.
	InternalErr = grpc.Errorf(codes.Internal, "server fault")

	// InvalidArgErr signals the caller that the arguments given with the request
	// are invalid.
	InvalidArgErr = grpc.Errorf(codes.InvalidArgument, "invalid arguments")

	// InvalidTokenErr signals the caller that its token is invalid, and a new
	// session should be started through an invocation of Login().
	InvalidTokenErr = grpc.Errorf(codes.PermissionDenied, "invalid token")

	// NotLoggedErr means that the caller tried to execute any operation different
	// from Login without a valid token.
	NotLoggedErr = grpc.Errorf(codes.Unauthenticated, "not authenticated")
)
