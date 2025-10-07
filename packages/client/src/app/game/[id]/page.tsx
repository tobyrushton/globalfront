import { FC } from "react";
import { GameWrapper } from "@/components/game/wrapper";
import { GameBoard } from "@/components/game/board";
import { serverClient } from "@/grpc/client";
import { GameProvider } from "@/components/game/provider";

type Props = {
    params: {
        id: string;
    }
}

const GamePage: FC<Props> = async ({ params: { id }}) => {
    const res = await serverClient.getGameDetails({ gameId: id })
    // TODO: handle res
    return (
        <GameProvider url={res.response.url}>
            <GameWrapper>
                <GameBoard />
            </GameWrapper>
        </GameProvider>
    )    
}

export default GamePage