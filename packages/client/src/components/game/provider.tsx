"use client"

import { FC, PropsWithChildren, createContext, useContext, useEffect } from "react"
import { WebsocketMessage } from "@globalfront/pb/messages/v1/messages"

type GameProviderProps = {
    url: string
}

type TGameContext = {}

const GameContext = createContext<TGameContext | null>(null)

export const GameProvider: FC<PropsWithChildren<GameProviderProps>> = ({ children, url }) => {
    useEffect(() => {
        console.log("Connecting to game server at", url)
        const socket = new WebSocket(url)
        socket.binaryType = "arraybuffer"

        socket.onopen = () => {
            console.log("Connected to game server")
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
