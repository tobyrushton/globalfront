package main

import (
	"log"
	"net"

	"github.com/tobyrushton/globalfront/packages/matchmaker"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
	"google.golang.org/grpc"
)

func main() {
	gf := gamefactory.New(60)
	mm := matchmaker.New(gf)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMatchmakerServer(grpcServer, mm)
	grpcServer.Serve(lis)
}
