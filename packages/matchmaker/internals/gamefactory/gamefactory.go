package gamefactory

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	pb "github.com/tobyrushton/globalfront/pb/game/v1"
)

type GameFactory struct {
	gameChannel  chan *pb.Game
	newGameChan  chan struct{}
	newGameDelay int
}

func New(newGameDelay int) *GameFactory {
	gf := &GameFactory{
		gameChannel:  make(chan *pb.Game),
		newGameChan:  make(chan struct{}),
		newGameDelay: newGameDelay,
	}
	go gf.createGameLoop()
	return gf
}

func (gf *GameFactory) createGame() {
	id := uuid.New().String()
	currentPlayers := int32(0)
	maxPlayers := int32(rand.IntN(32))
	game := &pb.Game{
		Id:          id,
		PlayerCount: currentPlayers,
		MaxPlayers:  maxPlayers,
	}

	gf.gameChannel <- game
}

func (gf *GameFactory) createGameLoop() {
	for {
		gf.createGame()

		timer := time.NewTimer(time.Duration(gf.newGameDelay) * time.Second)
		select {
		case <-timer.C:
			continue
		case <-gf.newGameChan:
			if !timer.Stop() {
				<-timer.C
			}
		}
	}
}

func (gf *GameFactory) GetNewGameChannel() chan<- struct{} {
	return gf.newGameChan
}

func (gf *GameFactory) GetGameChannel() <-chan *pb.Game {
	return gf.gameChannel
}
