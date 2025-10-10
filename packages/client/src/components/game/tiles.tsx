"use client"

import { FC } from "react"
import { useTiles } from "./tile-provider"

export const GameTiles: FC = () => {
    const { tiles } = useTiles()

    return (
        <div className="absolute">
            {tiles.map((row, rowIndex) => (
                <div key={rowIndex} className="flex">
                    {row.map((_, colIndex) => (
                        <div 
                            key={colIndex} 
                            className="w-[10px] h-[10px]"
                        />
                    ))}
                </div>
            ))}
        </div>
    )
}