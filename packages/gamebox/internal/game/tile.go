package game

type Tile struct {
	playerId string
	changed  bool
}

func NewTile() *Tile {
	return &Tile{}
}

func (t *Tile) PlayerId() string {
	return t.playerId
}

func (t *Tile) SetPlayerId(playerId string) {
	t.playerId = playerId
	t.changed = true
}

func (t *Tile) Changed() bool {
	return t.changed
}

func (t *Tile) ClearChanged() {
	t.changed = false
}
