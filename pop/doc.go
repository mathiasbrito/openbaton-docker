// Package pop implements a client-server protocol to remotely manage a docker daemon and use it as a Point of Presence for OpenBaton,
// abstracting Docker containers under OpenBaton catalogue types.
// pop reduces the attack surface of the of the infrastructure by acting as a proxy to a Docker daemon, implementing user based authentication
// and allowing clients to only operate on a small, proto-defined set of RPC invocations, and avoiding unnecessary exposure of the Docker daemon to the outside world.
// pop uses gRPC; see proto for the protocol definitions.
// This package is experimental. TLS is not ready yet, so it's not production ready (never use popd on a public network without TLS!).
package pop

//go:generate protoc -I ./proto ./proto/pop.proto --go_out=plugins=grpc:./proto
