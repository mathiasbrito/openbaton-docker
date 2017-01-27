// Package client implements a client for a pop service.
// Client wraps the gRPC client stub for pop, converting the container-oriented pop types into the OpenStack like structures OpenBaton
// expects.
// The clients are handled by a cache, to allow a stateless user like plugind to create clients on demand without incurring in additional
// overheads.
package client
