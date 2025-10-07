package gamemanager

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
	game "github.com/tobyrushton/globalfront/pb/game/v1"
	pb "github.com/tobyrushton/globalfront/pb/gamebox/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StartedGame struct {
	Game *game.Game
	Url  string
}

type GameManager struct {
	ctx context.Context

	gamesMu sync.Mutex
	games   map[string]StartedGame

	currentGame *game.Game

	playerMu sync.Mutex
	players  map[string]chan Update

	gf *gamefactory.GameFactory
}

type Update struct {
	GameId string
	Err    error
}

func NewGameManager(ctx context.Context, gf *gamefactory.GameFactory) *GameManager {
	gm := &GameManager{
		ctx:   ctx,
		gf:    gf,
		games: make(map[string]StartedGame),
	}

	go func() {
		for game := range gf.GetGameChannel() {
			if gm.currentGame != nil && gm.currentGame.PlayerCount > 1 {
				gm.startGame()
			}
			gm.gamesMu.Lock()
			if gm.currentGame != nil && gm.currentGame.PlayerCount == 1 {
				gm.playerMu.Lock()
				for _, ch := range gm.players {
					ch <- Update{
						GameId: "",
						Err:    errors.New("game cancelled due to insufficient players"),
					}
				}
				gm.playerMu.Unlock()
			}
			gm.currentGame = game
			gm.players = make(map[string]chan Update)
			gm.gamesMu.Unlock()
		}
	}()
	return gm
}

func (gm *GameManager) GetCurrentGame() *game.Game {
	return gm.currentGame
}

func (gm *GameManager) JoinGame() (string, error) {
	gm.gamesMu.Lock()

	if gm.currentGame == nil {
		return "", errors.New("no game available")
	}

	if len(gm.players) >= int(gm.currentGame.MaxPlayers) {
		return "", errors.New("game is full")
	}

	playerID := uuid.New().String()
	gm.playerMu.Lock()
	gm.players[playerID] = make(chan Update)
	gm.playerMu.Unlock()
	gm.currentGame.PlayerCount++

	gm.gamesMu.Unlock()
	if gm.currentGame.PlayerCount == gm.currentGame.MaxPlayers {
		gm.startGame()
	}

	return playerID, nil
}

func (gm *GameManager) startGame() error {
	fmt.Println("Attempting to start game...")
	gm.gamesMu.Lock()
	defer gm.gamesMu.Lock()
	fmt.Println("obtained lock to start game")

	conn, err := grpc.NewClient(
		"gamebox:5432",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewGameboxClient(conn)
	req := &pb.CreateGameRequest{
		Game: gm.currentGame,
	}
	res, err := client.CreateGame(gm.ctx, req)
	if err != nil {
		return err
	}

	gm.games[res.GameId] = StartedGame{
		Game: gm.currentGame,
		Url:  fmt.Sprintf("ws://gamebox:%d/ws", res.Port),
	}

	gm.playerMu.Lock()
	for _, ch := range gm.players {
		ch <- Update{
			GameId: res.GameId,
			Err:    nil,
		}
	}
	gm.playerMu.Unlock()

	return nil
}

func (gm *GameManager) RemovePlayer(playerId string) {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()
	gm.playerMu.Lock()
	defer gm.playerMu.Unlock()

	delete(gm.players, playerId)
	gm.currentGame.PlayerCount--
}

func (gm *GameManager) GetUpdateChannel(playerId string) chan Update {
	gm.playerMu.Lock()
	defer gm.playerMu.Unlock()

	return gm.players[playerId]
}

func (gm *GameManager) GetGame(gameId string) (StartedGame, bool) {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()

	game, exists := gm.games[gameId]
	return game, exists
}
