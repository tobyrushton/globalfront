package game

import (
	"fmt"
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
	http.HandleFunc("/ws", g.wsServer.ServeWS)

	return http.ListenAndServe(
		fmt.Sprintf(":%d", g.port),
		nil,
	)
}

func (g *Game) GetId() string {
	return g.game.Id
}
