"use client"

import { FC } from "react"
import { useTiles } from "./tile-provider"
import { usePlayers } from "./player-provider"
import { useGame } from "./provider"
import { useStatus } from "./status-provider"
import { MessageType, Spawn, WebsocketMessage } from "@globalfront/pb/messages/v1/messages"
import { convertCoordinatesToTileId } from "@/lib/tiles"

export const GameTiles: FC = () => {
    const { tiles } = useTiles()
    const { players } = usePlayers()
    const { send, playerId, player, attackPercentage } = useGame()
    const { gameStarted } = useStatus()

    const onClick = (row: number, col: number) => {
        if (!gameStarted) {
            send(WebsocketMessage.create({
                type: MessageType.MESSAGE_SPAWN,
                payload: {
                    oneofKind: "spawn",
                    spawn: Spawn.create({
                        playerId, tileId: convertCoordinatesToTileId(row, col)
                    })
                }
            }))
        } else {
            send(WebsocketMessage.create({
                type: MessageType.MESSAGE_ATTACK,
                payload: {
                    oneofKind: "attack",
                    attack: {
                        playerId, 
                        tileId: convertCoordinatesToTileId(row, col),
                        troopCount: Math.floor(player!.troopCount * (attackPercentage / 100))
                    }
                }
            }))
        }
    }

    return (
        <div className="absolute">
            {tiles.map((row, rowIndex) => (
                <div key={rowIndex} className="flex">
                    {row.map((_, colIndex) => (
                        <div 
                            key={colIndex} 
                            className="w-[10px] h-[10px]"
                            onClick={() => onClick(rowIndex, colIndex)}
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