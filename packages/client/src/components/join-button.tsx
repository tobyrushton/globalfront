"use client"

import { FC } from "react"
import { Game } from "@globalfront/pb/game/v1/game"
import { webClient } from "@/grpc/webClient"
import { useRouter } from "next/navigation"

type JoinButtonProps = {
    game: Game
}

export const JoinButton: FC<JoinButtonProps> = ({ game }) => {
    const router = useRouter()

    const joinGame = () => {
        const stream = webClient.joinGame({})

        stream.responses.onNext((response) => {
            if (response) {
                if (response.update.oneofKind === "serverDetails") 
                    router.push(`/game/${response.update.serverDetails.id}`)
            }
        })
    }

    return (
        <button 
            className='flex p-2 rounded-xl border-white cursor-pointer border-2'
            onClick={joinGame}
        >
            <div className='flex flex-col'>
            <p>game: {game.id}</p>
            <p>{game.playerCount} / {game.maxPlayers}</p>
            </div>
        </button>
    )
}