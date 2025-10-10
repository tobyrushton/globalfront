"use client"

import { FC, PropsWithChildren, createContext, useState, useContext } from "react"

type TTileContext = {
    tiles: string[][]
}

const TileContext = createContext<TTileContext | null>(null)

export const TileProvider: FC<PropsWithChildren> = ({ children }) => {
    const [tiles] = useState<string[][]>(
        Array.from({ length: 200 }, () => Array.from({ length: 200 }, () => ""))
    )

    return (
        <TileContext.Provider value={{ tiles }}>
            {children}
        </TileContext.Provider>
    )
}

export const useTiles = (): TTileContext => useContext(TileContext) as TTileContext
