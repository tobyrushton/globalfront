"use client"

import { FC, useEffect, useRef } from "react"
import { useTiles } from "./tile-provider"
import { usePlayers } from "./player-provider"
import { useGame } from "./provider"
import { useStatus } from "./status-provider"
import { MessageType, Spawn, WebsocketMessage } from "@globalfront/pb/messages/v1/messages"
import { convertCoordinatesToTileId } from "@/lib/tiles"

export const GameTiles: FC = () => {
    const canvasRef = useRef<HTMLCanvasElement>(null)

    const { tiles } = useTiles()
    const { players } = usePlayers()
    const { send, playerId, player, attackPercentage, scale } = useGame()
    const { gameStarted } = useStatus()

    const cellSize = 10;

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

    useEffect(() => {
        const canvas = canvasRef.current as HTMLCanvasElement
        const ctx = canvas.getContext("2d") as CanvasRenderingContext2D

        canvas.width = tiles[0].length * cellSize
        canvas.height = tiles.length * cellSize

        tiles.forEach((row, y) => {
            row.forEach((tile, x) => {
                if (tile != "0") {
                    ctx.fillStyle = players.get(tile)?.color || "#000000"
                    ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize)
                }
            })
        })

    }, [tiles])

    const handleClick = (e: MouseEvent) => {
        
    }

    return (
        <canvas 
            className="absolute"
            ref={canvasRef}
            onClick={e => {
                const rect = canvasRef.current?.getBoundingClientRect()
                const cssX = (e.clientX - (rect?.left || 0)) / scale
                const cssY = (e.clientY - (rect?.top || 0)) / scale
                const x = Math.floor(cssX / cellSize)
                const y = Math.floor(cssY / cellSize)

                onClick(y, x)
            }}
        />
    )
}