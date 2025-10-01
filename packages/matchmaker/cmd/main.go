package main

import (
	"log"
	"net"

	"github.com/tobyrushton/globalfront/packages/matchmaker"
	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
	"google.golang.org/grpc"
)

func main() {

	mm := matchmaker.New()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMatchmakerServer(grpcServer, mm)
	grpcServer.Serve(lis)
}
