"use client"

import { FC } from "react"
import { Game } from "@globalfront/pb/game/v1/game"
import { webClient } from "@/grpc/webClient"

type JoinButtonProps = {
    game: Game
}

const joinGame = () => {
    const stream = webClient.joinGame({})

    stream.responses.onNext((response) => {
        console.log(response)
    })
}

export const JoinButton: FC<JoinButtonProps> = ({ game }) => {
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