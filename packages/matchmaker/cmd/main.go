package main

import (
	"context"
	"log"
	"net"

	"github.com/tobyrushton/globalfront/packages/matchmaker"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamemanager"
	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
	"google.golang.org/grpc"
)

func main() {
	gf := gamefactory.New(60)
	gm := gamemanager.NewGameManager(context.Background(), gf)
	mm := matchmaker.New(gm)

	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMatchmakerServer(grpcServer, mm)
	grpcServer.Serve(lis)
}
