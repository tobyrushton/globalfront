package matchmaker

import (
	"context"

	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
)

var _ pb.MatchmakerServer = &MatchmakerServer{}

type MatchmakerServer struct {
	pb.UnimplementedMatchmakerServer
}

func New() *MatchmakerServer {
	return &MatchmakerServer{}
}

func (s *MatchmakerServer) GetCurrentGame(ctx context.Context, _ *pb.GetCurrentGameRequest) (*pb.GetCurrentGameResponse, error) {
	return &pb.GetCurrentGameResponse{
		Name: "Test Game",
	}, nil
}
