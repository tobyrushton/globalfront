"use client"

import { FC, MouseEvent, useState } from "react"
import { Game } from "@globalfront/pb/game/v1/game"
import { webClient } from "@/grpc/webClient"
import { useRouter } from "next/navigation"
import { cn } from "@/lib/utils"

type JoinButtonProps = {
    game: Game
}

export const JoinButton: FC<JoinButtonProps> = ({ game }) => {
    const [joined, setJoined] = useState(false)

    const router = useRouter()

    const joinGame = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        setJoined(true)
        const stream = webClient.joinGame({})

        stream.responses.onNext((response) => {
            if (response) {
                console.log(response)
                if (response.update.oneofKind === "serverDetails") {
                    console.log("pushing")
                    router.replace(`/game/${response.update.serverDetails.id}?playerId=${response.update.serverDetails.playerId}`)
                }
            }
        })

        // Add error handler
        stream.responses.onError((error) => {
            console.error("Stream error:", error)
        })
        
        // Add complete handler
        stream.responses.onComplete(() => {
            console.log("Stream completed")
        })
    }

    return (
        <button 
            className={
                cn(
                    'flex p-2 rounded-xl border-white cursor-pointer border-2', 
                    { 'bg-green-500': joined }
                )
            }
            onClick={joinGame}
        >
            <div className='flex flex-col'>
            <p>game: {game.id}</p>
            <p>{game.playerCount} / {game.maxPlayers}</p>
            </div>
        </button>
    )
}