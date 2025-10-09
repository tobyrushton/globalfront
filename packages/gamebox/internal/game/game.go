package game

import (
	"fmt"
	"net"
	"net/http"

	ws "github.com/tobyrushton/globalfront/packages/gamebox/internal/ws/server"
	pb "github.com/tobyrushton/globalfront/pb/game/v1"
)

type Game struct {
	port int
	game *pb.Game

	wsServer *ws.WsServer
}

func New(port int, game *pb.Game) *Game {
	return &Game{
		port:     port,
		game:     game,
		wsServer: ws.NewServer(),
	}
}

func (g *Game) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return err
	}
	s := &http.Server{
		Handler: g.wsServer,
	}

	return s.Serve(l)
}

func (g *Game) GetId() string {
	return g.game.Id
}

func (g *Game) GetPort() int {
	return g.port
}
