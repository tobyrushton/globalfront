"use client"

import { FC, PropsWithChildren, createContext, useContext, useState, useMemo, useEffect, useEffectEvent } from "react"
import { Player } from "@globalfront/pb/game/v1/game"

type TPlayerContext = {
    players: Map<string, Player>
    setPlayers: (players: Map<string, Player>) => void
    updatePlayerCounts: (updates: { [key: string]: number }) => void
}

const PlayerContext = createContext<TPlayerContext | null>(null)

export const PlayerProvider: FC<PropsWithChildren> = ({ children }) => {
    const [players, setPlayers] = useState<Map<string,Player>>(new Map<string, Player>())

    const updatePlayerCounts = (updates: { [key: string]: number }) => {
        if (Object.keys(updates).length === 0) return

        setPlayers(prev => {
            const tmpPlayers = new Map(prev)

            for (const [playerId, count] of Object.entries(updates)) {
                const player = tmpPlayers.get(playerId)
                if (player) {
                    tmpPlayers.set(playerId, { ...player, troopCount: count })
                }
            }

            return tmpPlayers
        })
    }

    const value = useMemo(() => ({ players, setPlayers, updatePlayerCounts }), [players])

    return (
        <PlayerContext.Provider value={value}>
            {children}
        </PlayerContext.Provider>
    )
}

export const usePlayers = (): TPlayerContext => useContext(PlayerContext) as TPlayerContext
