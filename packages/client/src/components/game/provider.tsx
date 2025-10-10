"use client"

import { FC, PropsWithChildren, createContext, useContext, useEffect } from "react"
import { WebsocketMessage, MessageType, JoinGame } from "@globalfront/pb/messages/v1/messages"

type GameProviderProps = {
    url: string
    playerId: string
}

type TGameContext = {}

const GameContext = createContext<TGameContext | null>(null)

export const GameProvider: FC<PropsWithChildren<GameProviderProps>> = ({ children, url, playerId }) => {
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
                console.log("Received message from server:", message)
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
