import { FC } from "react";
import { GameWrapper } from "@/components/game/wrapper";
import { GameBoard } from "@/components/game/board";
import { serverClient } from "@/grpc/client";
import { GameProvider } from "@/components/game/provider";

type Props = {
    params: Promise<{
        id: string;
    }>
    searchParams: Promise<{
        [key: string]: string | string[] | undefined
    }>
}

const GamePage: FC<Props> = async ({ params, searchParams }) => {
    const { id } = await params
    const sp = await searchParams
    const res = await serverClient.getGameDetails({ gameId: id })
    // TODO: handle res
    console.log("game url:",  res.response.url)
    return (
        <GameProvider url={res.response.url} playerId={sp?.playerId as string}>
            <GameWrapper>
                <GameBoard />
            </GameWrapper>
        </GameProvider>
    )    
}

export default GamePage