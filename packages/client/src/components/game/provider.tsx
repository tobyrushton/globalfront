"use client"

import { FC, PropsWithChildren, createContext, useContext, useEffect, useRef } from "react"
import { WebsocketMessage, MessageType, JoinGame, JoinGameResponse } from "@globalfront/pb/messages/v1/messages"
import { Player, Board } from "@globalfront/pb/game/v1/game"
import { useStatus } from "./status-provider"
import { usePlayers } from "./player-provider"
import { useTiles } from "./tile-provider"

type GameProviderProps = {
    url: string
    playerId: string
}

type TGameContext = {
    send: (msg: WebsocketMessage) => void
    playerId: string
}

const GameContext = createContext<TGameContext | null>(null)

export const GameProvider: FC<PropsWithChildren<GameProviderProps>> = ({ children, url, playerId }) => {
    const socketRef = useRef<WebSocket | null>(null) 
    const { setCountdown, startGame } = useStatus()
    const { setPlayers, updatePlayerCounts } = usePlayers()
    const { setBoard, handleTileUpdate } = useTiles()

    const send = (msg: WebsocketMessage) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(WebsocketMessage.toBinary(msg))
        } else {
            console.warn("WebSocket is not open. Unable to send message.")
        }
    }

    const handleJoinGameResponse = (res: JoinGameResponse) => {
        const playerMap = new Map<string, Player>()
        res.players.forEach(player => {
            playerMap.set(player.id, player)
        })
        setPlayers(playerMap)
        setBoard(res.board as Board)
    }

    const handleMessage = (msg: WebsocketMessage) => {
        switch(msg.payload.oneofKind) {
            case "startCountdown":
                setCountdown(msg.payload.startCountdown.countdownSeconds)
                break
            case "joinGameResponse":
                handleJoinGameResponse(msg.payload.joinGameResponse)
                break
            case "update":
                handleTileUpdate(msg.payload.update.updatedTiles)
                updatePlayerCounts(msg.payload.update.troopCountChanges)
                break
            case "gameStart":
                startGame()
                break
            default:
                console.log("Unhandled message:", msg)
        }
    }

    useEffect(() => {
        const socket = new WebSocket(url)
        socket.binaryType = "arraybuffer"
        socketRef.current = socket

        socket.onopen = () => {
            const msg = WebsocketMessage.create({
                type: MessageType.MESSAGE_JOIN_GAME,
                payload: {
                    oneofKind: "joinGame",
                    joinGame: JoinGame.create({ playerId })
                }
            })
            socket.send(WebsocketMessage.toBinary(msg))
        }

        socket.onmessage = (event) => {
            if (event.data instanceof ArrayBuffer) {
                const message = WebsocketMessage.fromBinary(new Uint8Array(event.data))
                handleMessage(message)
            }
        }

        return () => socket.close()
    }, [url])
    
    return (
        <GameContext.Provider value={{ send, playerId }}>
            {children}
        </GameContext.Provider>
    )
}

export const useGame = (): TGameContext => useContext(GameContext) as TGameContext
