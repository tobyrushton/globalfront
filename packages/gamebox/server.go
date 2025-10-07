package gamebox

import (
	"context"
	"fmt"

	"github.com/tobyrushton/globalfront/packages/gamebox/internal/spawner"
	pb "github.com/tobyrushton/globalfront/pb/gamebox/v1"
)

type GameboxServer struct {
	pb.UnimplementedGameboxServer

	s *spawner.Spawner
}

var _ pb.GameboxServer = &GameboxServer{}

func New() *GameboxServer {
	s := spawner.New()
	return &GameboxServer{
		s: s,
	}
}

func (gb *GameboxServer) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	fmt.Println("CreateGame called with request:", req)
	game, err := gb.s.Spawn(req.Game)
	if err != nil {
		return nil, err
	}

	return &pb.CreateGameResponse{
		GameId: game.GetId(),
		Port:   int32(game.GetPort()),
	}, nil
}
