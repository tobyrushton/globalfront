package game

type Tile struct {
	playerId string
}

func NewTile() *Tile {
	return &Tile{}
}

func (t *Tile) PlayerId() string {
	return t.playerId
}

func (t *Tile) SetPlayerId(playerId string) {
	t.playerId = playerId
}
