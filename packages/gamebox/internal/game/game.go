package game

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	ws "github.com/tobyrushton/globalfront/packages/gamebox/internal/ws"
	pb "github.com/tobyrushton/globalfront/pb/game/v1"
	v1 "github.com/tobyrushton/globalfront/pb/messages/v1"
)

type Game struct {
	port int
	game *pb.Game

	wsServer *ws.WsServer

	started bool

	msgChan chan *v1.WebsocketMessage

	playersMu sync.Mutex
	players   map[string]*ws.Client
}

func New(port int, game *pb.Game, players []string) *Game {
	playerMap := make(map[string]*ws.Client)

	for _, player := range players {
		playerMap[player] = nil
	}

	msgChan := make(chan *v1.WebsocketMessage, 100)

	return &Game{
		port:     port,
		game:     game,
		wsServer: ws.NewServer(msgChan),
		started:  false,
		players:  playerMap,
		msgChan:  msgChan,
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

	go func() {
		for {
			msg := <-g.msgChan
			fmt.Println("Received message:", msg)
		}
	}()

	go g.startGame()

	return s.Serve(l)
}

func (g *Game) GetId() string {
	return g.game.Id
}

func (g *Game) GetPort() int {
	return g.port
}

func (g *Game) startGame() {
	for i := 60; i >= 1; i-- {
		fmt.Println("Starting game in", i, "seconds")
		g.wsServer.Broadcast(&v1.WebsocketMessage{
			Type: v1.MessageType_MESSAGE_START_COUNTDOWN,
			Payload: &v1.WebsocketMessage_StartCountdown{
				StartCountdown: &v1.StartCountdown{
					CountdownSeconds: int32(i),
				},
			},
		})
		time.Sleep(1 * time.Second)
	}
	g.wsServer.Broadcast(&v1.WebsocketMessage{
		Type:    v1.MessageType_MESSAGE_GAME_START,
		Payload: &v1.WebsocketMessage_GameStart{},
	})
}
