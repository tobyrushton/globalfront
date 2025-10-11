package game

import (
	"math"
	"sync"

	v1 "github.com/tobyrushton/globalfront/pb/game/v1"
)

const (
	spawnRadius = 5
)

type Board struct {
	tilesMu sync.Mutex
	tiles   [][]*Tile

	height int
	width  int
}

func NewBoard() *Board {
	tiles := make([][]*Tile, 200)
	for i := range tiles {
		tiles[i] = make([]*Tile, 200)
		for j := range tiles[i] {
			tiles[i][j] = NewTile()
		}
	}
	return &Board{
		tiles:  tiles,
		height: 200,
		width:  200,
	}
}

func (b *Board) Tiles() [][]*Tile {
	return b.tiles
}

func (b *Board) Board() *v1.Board {
	board := &v1.Board{
		Rows: make([]*v1.BoardRow, len(b.tiles)),
	}

	for i, row := range b.tiles {
		boardRow := &v1.BoardRow{
			Tiles: make([]*v1.Tile, len(row)),
		}
		for j, tile := range row {
			boardRow.Tiles[j] = &v1.Tile{
				PlayerId: tile.PlayerId(),
			}
		}
		board.Rows[i] = boardRow
	}

	return board
}

func (b *Board) SetPlayerSpawn(playerId string, tileId int32) {
	b.tilesMu.Lock()
	defer b.tilesMu.Unlock()

	centerX := int(tileId / 200)
	centerY := int(tileId % 200)

	for dy := -spawnRadius; dy <= spawnRadius; dy++ {
		for dx := -spawnRadius; dx <= spawnRadius; dx++ {
			x := centerX + dx
			y := centerY + dy

			// Skip out-of-bounds tiles
			if x < 0 || y < 0 || x >= b.width || y >= b.height {
				continue
			}

			// Check if the tile is within the circular radius
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			if distance <= float64(spawnRadius) {
				b.tiles[x][y].SetPlayerId(playerId)
			}
		}
	}
}
