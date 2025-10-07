"use client"

import { FC, PropsWithChildren, createContext, useContext, useEffect } from "react"
import useWebSocket from "react-use-websocket"

type GameProviderProps = {
    url: string
}

type TGameContext = {}

const GameContext = createContext<TGameContext | null>(null)

export const GameProvider: FC<PropsWithChildren<GameProviderProps>> = ({ children, url }) => {
    const { readyState } = useWebSocket(url, {
        share: false,
        shouldReconnect: () => true,
    })

    useEffect(() => {
        console.log("WebSocket readyState:", readyState)
    }, [readyState])
    
    return (
        <GameContext.Provider value={null}>
            {children}
        </GameContext.Provider>
    )
}

export const useGame = (): TGameContext => useContext(GameContext) as TGameContext
