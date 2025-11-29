package game

import (
	"math"
	"sync"

	v1 "github.com/tobyrushton/globalfront/pb/game/v1"
)

const (
	spawnRadius = 4
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

func (b *Board) clearPlayer(playerId string) {
	b.tilesMu.Lock()
	defer b.tilesMu.Unlock()

	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			if b.tiles[i][j].PlayerId() == playerId {
				b.tiles[i][j].SetPlayerId("")
			}
		}
	}
}

func (b *Board) SetPlayerSpawn(playerId string, tileId int32) {
	b.clearPlayer(playerId)
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

func (b *Board) GetChangedTiles() map[int32]string {
	b.tilesMu.Lock()
	defer b.tilesMu.Unlock()

	changedTiles := make(map[int32]string)

	for i, row := range b.tiles {
		for j, tile := range row {
			if tile.Changed() {
				tileId := int32(i*200 + j)
				changedTiles[tileId] = tile.PlayerId()
				tile.ClearChanged()
			}
		}
	}

	return changedTiles
}

func (b *Board) GetTile(tileId int32) *Tile {
	b.tilesMu.Lock()
	defer b.tilesMu.Unlock()

	x, y := CoordinatesFromTileId(tileId)

	return b.tiles[x][y]
}

func (b *Board) FindBorder(player1, player2 string, start int32) []int32 {
	height, width := len(b.tiles), len(b.tiles[0])

	borderTiles := make([]int32, 0)

	// complete a dfs in order to find the border cells
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	visited := make(map[int32]struct{})
	x, y := CoordinatesFromTileId(start)
	s := [][2]int{{x, y}}

	for len(s) > 0 {
		tile := s[len(s)-1]
		s = s[:len(s)-1]

		for _, dir := range dirs {
			nx, ny := tile[0]+dir[0], tile[1]+dir[1]
			nid := TileId(nx, ny)

			if _, seen := visited[nid]; !seen &&
				nx >= 0 &&
				ny >= 0 &&
				nx < height &&
				ny < width {
				visited[nid] = struct{}{}
				ntile := b.tiles[nx][ny]
				if ntile.PlayerId() == player2 {
					borderTiles = append(borderTiles, nid)
				}
			}
		}
	}

	return borderTiles
}

func TileId(x, y int) int32 {
	return int32(x*200 + y)
}

func CoordinatesFromTileId(tileId int32) (int, int) {
	x := int(tileId / 200)
	y := int(tileId % 200)
	return x, y
}
