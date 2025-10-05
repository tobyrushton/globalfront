package gamemanager

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
	game "github.com/tobyrushton/globalfront/pb/game/v1"
)

type GameManager struct {
	gameMu sync.Mutex
	game   *game.Game

	playerDetails map[string]struct{}

	gf *gamefactory.GameFactory
}

func NewGameManager(gf *gamefactory.GameFactory) *GameManager {
	gm := &GameManager{
		gf: gf,
	}

	go func() {
		for game := range gf.GetGameChannel() {
			gm.gameMu.Lock()
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
	defer gm.gameMu.Unlock()

	if gm.game == nil {
		return "", errors.New("no game available")
	}

	if len(gm.playerDetails) >= int(gm.game.MaxPlayers) {
		return "", errors.New("game is full")
	}

	playerID := uuid.New().String()
	gm.playerDetails[playerID] = struct{}{}
	gm.game.PlayerCount++

	return playerID, nil
}
