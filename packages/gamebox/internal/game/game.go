package game

import (
	"context"
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
	ctx context.Context

	port int
	game *pb.Game

	wsServer *ws.WsServer

	started bool

	msgChan chan *v1.WebsocketMessage

	playersMu sync.Mutex
	players   map[string]*pb.Player

	board *Board

	am *AttackManager
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

	board := NewBoard()
	am := NewAttackManager(board)

	return &Game{
		ctx:      context.TODO(),
		port:     port,
		game:     game,
		wsServer: ws.NewServer(msgChan, players),
		started:  false,
		players:  playerMap,
		msgChan:  msgChan,
		board:    board,
		am:       am,
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
	go g.updateLoop()

	return s.Serve(l)
}

func (g *Game) GetId() string {
	return g.game.Id
}

func (g *Game) GetPort() int {
	return g.port
}

func (g *Game) startGame() {
	for i := 30; i >= 1; i-- {
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
	g.started = true
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
	case *v1.WebsocketMessage_Attack:
		g.handleAttack(p.Attack.PlayerId, p.Attack.TileId, p.Attack.TroopCount)
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

func (g *Game) handleAttack(attackerId string, tileId int32, troopCount int32) {
	g.am.InitAttack(attackerId, tileId, troopCount)
}

func (g *Game) updateLoop() {
	ticker := time.NewTicker(time.Second / 20)
	tick := 0

	for {
		tick++
		select {
		case <-g.ctx.Done():
			return
		case <-ticker.C:
			msg := &v1.WebsocketMessage{
				Type: v1.MessageType_MESSAGE_UPDATE,
				Payload: &v1.WebsocketMessage_Update{
					Update: &v1.Update{},
				},
			}
			send := false
			if boardUpdates := g.board.GetChangedTiles(); len(boardUpdates) > 0 {
				msg.GetUpdate().UpdatedTiles = boardUpdates
				send = true
			}
			if g.started && tick%20 == 0 {
				if troopUpdates := g.calculateTroopUpdates(); len(troopUpdates) > 0 {
					msg.GetUpdate().TroopCountChanges = troopUpdates
					send = true
				}
			}
			if send {
				g.wsServer.Broadcast(msg)
			}
		}
	}
}

func (g *Game) calculateTroopUpdates() map[string]int32 {
	updates := make(map[string]int32)

	g.playersMu.Lock()
	defer g.playersMu.Unlock()

	for playerId := range g.players {
		updatedCount := int32(float32(g.players[playerId].TroopCount) * 1.05)
		updates[playerId] = updatedCount
		g.players[playerId].TroopCount = updatedCount
	}

	return updates
}
