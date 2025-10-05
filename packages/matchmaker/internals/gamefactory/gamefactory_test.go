package gamefactory_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tobyrushton/globalfront/packages/matchmaker/internals/gamefactory"
)

func TestCreateGamesOnTimer(t *testing.T) {
	gf := gamefactory.New(1)
	gameChan := gf.GetGameChannel()
	firstGame := <-gameChan

	require.NotNil(t, firstGame)
	secondGame := <-gameChan
	require.NotNil(t, secondGame)
	require.NotEqual(t, firstGame.Id, secondGame.Id)
}

func TestCreateGameOnRequest(t *testing.T) {
	gf := gamefactory.New(60)
	gameChan := gf.GetGameChannel()
	requestChan := gf.GetNewGameChannel()

	firstGame := <-gameChan
	require.NotNil(t, firstGame)

	timeRequested := time.Now()
	requestChan <- struct{}{}
	secondGame := <-gameChan
	require.NotNil(t, secondGame)
	require.NotEqual(t, firstGame.Id, secondGame.Id)
	timeTaken := time.Since(timeRequested)
	require.Less(t, timeTaken.Milliseconds(), int64(1000))
}
