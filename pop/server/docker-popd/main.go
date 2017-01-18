package main

import (
    "net"
    
	"google.golang.org/grpc"

    log "github.com/sirupsen/logrus"
)

func main() {
    lis, err := net.Listen("tcp", ":60000")
    if err != nil {
        log.Fatal(err)
    }

    srv := grpc.NewServer(
        grpc.StreamInterceptor(streamInterceptor),
        grpc.UnaryInterceptor(unaryInterceptor),
    )

    pop.RegisterPoPServer(srv, routes{})
    
    if err := srv.Serve(lis); err != nil {
        log.Fatal(err)
    }
}