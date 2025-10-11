package game

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	ws "github.com/tobyrushton/globalfront/packages/gamebox/internal/ws"
	"github.com/tobyrushton/globalfront/packages/gamebox/utils"
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
	players   map[string]*pb.Player

	board *Board
}

func New(port int, game *pb.Game, players []string) *Game {
	playerMap := make(map[string]*pb.Player)

	for _, player := range players {
		playerMap[player] = &pb.Player{
			Id:         player,
			Color:      utils.RandomColor(),
			TroopCount: 3000,
		}
	}

	msgChan := make(chan *v1.WebsocketMessage, 100)

	return &Game{
		port:     port,
		game:     game,
		wsServer: ws.NewServer(msgChan, players),
		started:  false,
		players:  playerMap,
		msgChan:  msgChan,
		board:    NewBoard(),
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
			g.handleMsg(msg)
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

func (g *Game) handleMsg(msg *v1.WebsocketMessage) {
	switch p := msg.Payload.(type) {
	case *v1.WebsocketMessage_JoinGame:
		err := g.joinPlayer(p.JoinGame.PlayerId)
		if err != nil {
			fmt.Println("Error joining player:", err)
		}
	case *v1.WebsocketMessage_Spawn:
		g.handleSpawn(p.Spawn.PlayerId, p.Spawn.TileId)
	default:
		fmt.Println("Unhandled message type:", msg.Type)
	}
}

func (g *Game) joinPlayer(playerId string) error {
	g.playersMu.Lock()
	defer g.playersMu.Unlock()

	if _, ok := g.players[playerId]; !ok {
		return fmt.Errorf("player %s not in game", playerId)
	}

	msg := &v1.WebsocketMessage{
		Type: v1.MessageType_MESSAGE_JOIN_GAME_RESPONSE,
		Payload: &v1.WebsocketMessage_JoinGameResponse{
			JoinGameResponse: &v1.JoinGameResponse{
				Players: utils.FlattenMap(g.players),
				Board:   g.board.Board(),
			},
		},
	}
	g.wsServer.SendToPlayer(playerId, msg)

	return nil
}

func (g *Game) handleSpawn(playerId string, tileId int32) {
	g.board.SetPlayerSpawn(playerId, tileId)
}
