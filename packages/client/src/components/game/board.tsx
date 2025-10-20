import { FC } from "react"
import { GameTiles } from "./tiles"

export const GameBoard: FC = () => {
    return (
        <div className="w-[2000px] h-[2000px] bg-green-500/20 border border-green-500 relative">
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-20 h-20 bg-green-500/50 border border-green-500 flex justify-center items-center">
                <span className="text-white font-bold">Center</span>
            </div>
            <GameTiles />
        </div>
    )
}