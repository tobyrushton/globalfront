package game

import (
	pb "github.com/tobyrushton/globalfront/pb/game/v1"
)

type Game struct {
	port int
	game *pb.Game
}

func New(port int, game *pb.Game) *Game {
	return &Game{
		port: port,
		game: game,
	}
}

func (g *Game) Start() error {
	return nil
}

func (g *Game) GetId() string {
	return g.game.Id
}
