"use client"

import { FC, PropsWithChildren, createContext, useContext, useEffect } from "react"
import { WebsocketMessage, MessageType, JoinGame } from "@globalfront/pb/messages/v1/messages"
import { Player } from "@globalfront/pb/game/v1/game"
import { useStatus } from "./status-provider"
import { usePlayers } from "./player-provider"

type GameProviderProps = {
    url: string
    playerId: string
}

type TGameContext = {}

const GameContext = createContext<TGameContext | null>(null)

export const GameProvider: FC<PropsWithChildren<GameProviderProps>> = ({ children, url, playerId }) => {
    const { setCountdown } = useStatus()
    const { setPlayers } = usePlayers()

    const handleJoinGameResponse = (players: Player[]) => {
        const playerMap = new Map<string, Player>()
        players.forEach(player => {
            playerMap.set(player.id, player)
        })
        setPlayers(playerMap)
    }

    const handleMessage = (msg: WebsocketMessage) => {
        switch(msg.payload.oneofKind) {
            case "startCountdown":
                setCountdown(msg.payload.startCountdown.countdownSeconds)
                break
            case "joinGameResponse":
                handleJoinGameResponse(msg.payload.joinGameResponse.players)
            default:
                console.log("Unhandled message:", msg)
        }
    }

    useEffect(() => {
        const socket = new WebSocket(url)
        socket.binaryType = "arraybuffer"

        socket.onopen = () => {
            const msg = WebsocketMessage.create({
                type: MessageType.MESSAGE_JOIN_GAME,
                payload: {
                    oneofKind: "joinGame",
                    joinGame: JoinGame.create({ playerId })
                }
            })
            console.log("Sending join game message:", msg)
            socket.send(WebsocketMessage.toBinary(msg))
        }

        socket.onmessage = (event) => {
            console.log(event)
            if (event.data instanceof ArrayBuffer) {
                const message = WebsocketMessage.fromBinary(new Uint8Array(event.data))
                handleMessage(message)
            }
        }

        return () => socket.close()
    }, [url])
    
    return (
        <GameContext.Provider value={null}>
            {children}
        </GameContext.Provider>
    )
}

export const useGame = (): TGameContext => useContext(GameContext) as TGameContext
