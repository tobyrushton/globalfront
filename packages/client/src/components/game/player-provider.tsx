"use client"

import { FC, PropsWithChildren, createContext, useContext, useState, useMemo } from "react"
import { Player } from "@globalfront/pb/game/v1/game"

type TPlayerContext = {
    players: Map<string, Player>
    setPlayers: (players: Map<string, Player>) => void
}

const PlayerContext = createContext<TPlayerContext | null>(null)

export const PlayerProvider: FC<PropsWithChildren> = ({ children }) => {
    const [players, setPlayers] = useState<Map<string,Player>>(new Map<string, Player>())

    const value = useMemo(() => ({ players, setPlayers }), [players])
    return (
        <PlayerContext.Provider value={value}>
            {children}
        </PlayerContext.Provider>
    )
}

export const usePlayers = (): TPlayerContext => useContext(PlayerContext) as TPlayerContext
