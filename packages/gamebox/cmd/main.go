package main

import (
	"log"
	"net"

	"github.com/tobyrushton/globalfront/packages/gamebox"
	pb "github.com/tobyrushton/globalfront/pb/gamebox/v1"
	"google.golang.org/grpc"
)

func main() {
	gb := gamebox.New()

	lis, err := net.Listen("tcp", ":5432")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGameboxServer(grpcServer, gb)
	grpcServer.Serve(lis)
}
