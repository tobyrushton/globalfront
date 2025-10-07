package spawner

import (
	"errors"
	"sync"

	"github.com/tobyrushton/globalfront/packages/gamebox/internal/game"
	pb "github.com/tobyrushton/globalfront/pb/game/v1"
)

type Spawner struct {
	gamesMu sync.Mutex
	games   map[string]*game.Game

	portsMu sync.Mutex
	ports   map[int]struct{}
}

func New() *Spawner {
	return &Spawner{
		games: make(map[string]*game.Game),
		ports: make(map[int]struct{}),
	}
}

func (s *Spawner) Spawn(gm *pb.Game) (*game.Game, error) {
	s.gamesMu.Lock()
	defer s.gamesMu.Unlock()

	if _, exists := s.games[gm.Id]; exists {
		return nil, nil
	}

	port := s.findFreePort()
	if port == 0 {
		return nil, errors.New("no free ports available")
	}

	newGame := game.New(port, gm)
	go newGame.Start()

	s.games[gm.Id] = newGame
	return newGame, nil
}

func (s *Spawner) findFreePort() int {
	s.portsMu.Lock()
	defer s.portsMu.Unlock()

	for port := 5433; port <= 5533; port++ {
		if _, inUse := s.ports[port]; !inUse {
			s.ports[port] = struct{}{}
			return port
		}
	}
	return 0 // No free port found
}
