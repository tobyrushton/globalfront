"use client"

import { FC, PropsWithChildren, useState, createContext, useContext, useMemo } from "react"

type TStatusContext = {
    setCountdown: (seconds: number) => void
    startGame: () => void
    endGame: () => void
    gameStarted: boolean
}

const StatusContext = createContext<TStatusContext | null>(null)

export const StatusProvider: FC<PropsWithChildren> = ({ children }) => {
    const [countdown, setCountdown] = useState<number>(60)
    const [gameStarted, setGameStarted] = useState<boolean>(false)
    
    const startGame = () => setGameStarted(true)
    const endGame = () => {}
   
    const value = useMemo(() => ({
        setCountdown,
        startGame,
        endGame,
        gameStarted
    }), [gameStarted])

    return (
        <StatusContext.Provider value={value}>
            <div 
                className="absolute h-2 bg-sky-500"
                style={{
                    "width": `${((60 - countdown )/ 60) * 100}%`
                }}
            />
            {children}
        </StatusContext.Provider>
    )
}

export const useStatus = (): TStatusContext => useContext(StatusContext) as TStatusContext
