package matchmaker

import (
	"context"

	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
	game "github.com/tobyrushton/globalfront/pb/game/v1"
	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
)

var _ pb.MatchmakerServer = &MatchmakerServer{}

type MatchmakerServer struct {
	pb.UnimplementedMatchmakerServer

	gf   *gamefactory.GameFactory
	game *game.Game
}

func New(gf *gamefactory.GameFactory) *MatchmakerServer {
	mm := &MatchmakerServer{
		gf: gf,
	}
	go func() {
		for game := range gf.GetGameChannel() {
			mm.game = game
		}
	}()

	return mm
}

func (s *MatchmakerServer) GetCurrentGame(ctx context.Context, _ *pb.GetCurrentGameRequest) (*pb.GetCurrentGameResponse, error) {
	return &pb.GetCurrentGameResponse{
		Game: s.game,
	}, nil
}
