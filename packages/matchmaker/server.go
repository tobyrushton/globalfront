package matchmaker

import (
	"context"

	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamemanager"
	pb "github.com/tobyrushton/globalfront/pb/matchmaker/v1"
)

var _ pb.MatchmakerServer = &MatchmakerServer{}

type MatchmakerServer struct {
	pb.UnimplementedMatchmakerServer

	gm *gamemanager.GameManager
}

func New(gm *gamemanager.GameManager) *MatchmakerServer {
	return &MatchmakerServer{
		gm: gm,
	}
}

func (s *MatchmakerServer) GetCurrentGame(ctx context.Context, _ *pb.GetCurrentGameRequest) (*pb.GetCurrentGameResponse, error) {
	return &pb.GetCurrentGameResponse{
		Game: s.gm.GetCurrentGame(),
	}, nil
}

func (s *MatchmakerServer) JoinGame(req *pb.JoinGameRequest, stream pb.Matchmaker_JoinGameServer) error {
	playerID, err := s.gm.JoinGame()
	joinUpdate := &pb.JoinUpdate{}
	isErrJoining := false
	if err != nil {
		isErrJoining = true
		joinUpdate.Update = &pb.JoinUpdate_Error{
			Error: &pb.JoinError{
				Message: err.Error(),
			},
		}
	} else {
		joinUpdate.Update = &pb.JoinUpdate_Acknowledgement{
			Acknowledgement: &pb.JoinAcknowledgement{
				Message: "Joined game successfully",
			},
		}
	}

	if err := stream.Send(joinUpdate); err != nil {
		return err
	}

	if isErrJoining {
		return nil
	}

	select {
	case <-stream.Context().Done():
		s.gm.RemovePlayer(playerID)
		return stream.Context().Err()
	case update := <-s.gm.GetUpdateChannel(playerID):
		joinUpdate = &pb.JoinUpdate{}
		if update.Err != nil {
			joinUpdate.Update = &pb.JoinUpdate_Error{
				Error: &pb.JoinError{
					Message: update.Err.Error(),
				},
			}
		} else {
			joinUpdate.Update = &pb.JoinUpdate_ServerDetails{
				ServerDetails: &pb.ServerDetails{
					Id: update.GameId,
				},
			}
		}
		if err := stream.Send(joinUpdate); err != nil {
			return err
		}
	}

	return nil
}
