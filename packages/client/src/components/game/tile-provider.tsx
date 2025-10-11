"use client"

import { FC, PropsWithChildren, createContext, useState, useContext, useMemo } from "react"
import { Board } from "@globalfront/pb/game/v1/game"

type TTileContext = {
    tiles: string[][]
    setBoard: (board: Board) => void
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
    
    const value = useMemo(() => ({ tiles, setBoard }), [tiles])

    return (
        <TileContext.Provider value={value}>
            {children}
        </TileContext.Provider>
    )
}

export const useTiles = (): TTileContext => useContext(TileContext) as TTileContext
