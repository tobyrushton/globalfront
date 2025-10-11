"use client"

import { FC } from "react"
import { useTiles } from "./tile-provider"
import { usePlayers } from "./player-provider"

export const GameTiles: FC = () => {
    const { tiles } = useTiles()
    const { players } = usePlayers()

    return (
        <div className="absolute">
            {tiles.map((row, rowIndex) => (
                <div key={rowIndex} className="flex">
                    {row.map((_, colIndex) => (
                        <div 
                            key={colIndex} 
                            className="w-[10px] h-[10px]"
                            style={{ 
                                backgroundColor: players.get(tiles[rowIndex][colIndex])?.color || 'transparent' 
                            }}
                        />
                    ))}
                </div>
            ))}
        </div>
    )
}