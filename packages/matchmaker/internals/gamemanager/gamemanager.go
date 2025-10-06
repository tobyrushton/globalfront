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

	playerDetails map[string]struct{}

	gf *gamefactory.GameFactory

	updateChan chan Update
}

type Update struct {
	GameId string
	Err    error
}

func NewGameManager(ctx context.Context, gf *gamefactory.GameFactory) *GameManager {
	gm := &GameManager{
		ctx:        ctx,
		gf:         gf,
		updateChan: make(chan Update, 1),
	}

	go func() {
		for game := range gf.GetGameChannel() {
			if gm.game != nil && gm.game.PlayerCount > 1 {
				gm.startGame()
			}
			gm.gameMu.Lock()
			if gm.game != nil && gm.game.PlayerCount == 1 {
				gm.updateChan <- Update{
					GameId: "",
					Err:    errors.New("not enough players to start a game"),
				}
			}
			gm.game = game
			gm.playerDetails = make(map[string]struct{})
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

	if len(gm.playerDetails) >= int(gm.game.MaxPlayers) {
		return "", errors.New("game is full")
	}

	playerID := uuid.New().String()
	gm.playerDetails[playerID] = struct{}{}
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

	gm.updateChan <- Update{
		GameId: res.GameId,
		Err:    nil,
	}

	return nil
}

func (gm *GameManager) GetUpdateChannel() <-chan Update {
	return gm.updateChan
}

func (gm *GameManager) RemovePlayer(playerId string) {
	gm.gameMu.Lock()
	defer gm.gameMu.Unlock()

	delete(gm.playerDetails, playerId)
	gm.game.PlayerCount--
}
