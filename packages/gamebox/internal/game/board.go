package game

import v1 "github.com/tobyrushton/globalfront/pb/game/v1"

type Board struct {
	tiles [][]*Tile
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
		tiles: tiles,
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
