"use client"

import { FC, PropsWithChildren, createContext, useState, useContext, useMemo } from "react"
import { Board } from "@globalfront/pb/game/v1/game"

type TTileContext = {
    tiles: string[][]
    setBoard: (board: Board) => void
    handleTileUpdate: (updates: { [key: number]: string }) => void
}

const TileContext = createContext<TTileContext | null>(null)

export const TileProvider: FC<PropsWithChildren> = ({ children }) => {
    const [tiles, setTiles] = useState<string[][]>(
        Array.from({ length: 200 }, () => Array.from({ length: 200 }, () => ""))
    )

    const setBoard = (board: Board) => {
        const newTiles = Array.from({ length: 200 }, () => Array.from({ length: 200 }, () => ""))
        board.rows.forEach((row, i) => {
            row.tiles.forEach((tile, j) => {
                newTiles[i][j] = tile.playerId
            })
        })
        setTiles(newTiles)
    }

    const handleTileUpdate = (updates: { [key: number]: string }) => {
        setTiles(prevTiles => {
            const newTiles = prevTiles.map(row => [...row])

            for (const [tileIdStr, playerId] of Object.entries(updates)) {
                const tileId = parseInt(tileIdStr, 10)
                const row = Math.floor(tileId / 200)
                const col = tileId % 200
                if (row >= 0 && row < 200 && col >= 0 && col < 200) {
                    newTiles[row][col] = playerId
                }
            }

            return newTiles
        })
    }
    
    const value = useMemo(() => ({ tiles, setBoard, handleTileUpdate }), [tiles])

    return (
        <TileContext.Provider value={value}>
            {children}
        </TileContext.Provider>
    )
}

export const useTiles = (): TTileContext => useContext(TileContext) as TTileContext
