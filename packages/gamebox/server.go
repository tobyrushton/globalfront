package gamebox

import (
	"context"

	pb "github.com/tobyrushton/globalfront/pb/gamebox/v1"
)

type GameboxServer struct {
	pb.UnimplementedGameboxServer
}

var _ pb.GameboxServer = &GameboxServer{}

func New() *GameboxServer {
	return &GameboxServer{}
}

func (s *GameboxServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	return &pb.CreateGameResponse{
		GameId: "123",
	}, nil
}
