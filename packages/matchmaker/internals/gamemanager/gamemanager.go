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

type GameManager struct {
	ctx context.Context

	gameMu sync.Mutex
	game   *game.Game

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
		ctx: ctx,
		gf:  gf,
	}

	go func() {
		for game := range gf.GetGameChannel() {
			if gm.game != nil && gm.game.PlayerCount > 1 {
				gm.startGame()
			}
			gm.gameMu.Lock()
			if gm.game != nil && gm.game.PlayerCount == 1 {
				gm.playerMu.Lock()
				for _, ch := range gm.players {
					ch <- Update{
						GameId: "",
						Err:    errors.New("game cancelled due to insufficient players"),
					}
				}
				gm.playerMu.Unlock()
			}
			gm.game = game
			gm.players = make(map[string]chan Update)
			gm.gameMu.Unlock()
		}
	}()
	return gm
}

func (gm *GameManager) GetCurrentGame() *game.Game {
	return gm.game
}

func (gm *GameManager) JoinGame() (string, error) {
	gm.gameMu.Lock()

	if gm.game == nil {
		return "", errors.New("no game available")
	}

	if len(gm.players) >= int(gm.game.MaxPlayers) {
		return "", errors.New("game is full")
	}

	playerID := uuid.New().String()
	gm.playerMu.Lock()
	gm.players[playerID] = make(chan Update)
	gm.playerMu.Unlock()
	gm.game.PlayerCount++

	gm.gameMu.Unlock()
	if gm.game.PlayerCount == gm.game.MaxPlayers {
		gm.startGame()
	}

	return playerID, nil
}

func (gm *GameManager) startGame() error {
	fmt.Println("Attempting to start game...")
	gm.gameMu.Lock()
	defer gm.gameMu.Lock()
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
		Game: gm.game,
	}
	fmt.Println("Starting game with ID:", gm.game.Id)
	res, err := client.CreateGame(gm.ctx, req)
	if err != nil {
		fmt.Println("Error starting game:", err)
		return err
	}
	fmt.Println("Game started with ID:", res.GameId)

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
	gm.gameMu.Lock()
	defer gm.gameMu.Unlock()
	gm.playerMu.Lock()
	defer gm.playerMu.Unlock()

	delete(gm.players, playerId)
	gm.game.PlayerCount--
}

func (gm *GameManager) GetUpdateChannel(playerId string) chan Update {
	gm.gameMu.Lock()
	defer gm.gameMu.Unlock()

	return gm.players[playerId]
}
