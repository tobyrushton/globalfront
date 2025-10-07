import { FC } from "react";
import { GameWrapper } from "@/components/game/wrapper";
import { GameBoard } from "@/components/game/board";
import { serverClient } from "@/grpc/client";

type Props = {
    params: {
        id: string;
    }
}

const GamePage: FC<Props> = async ({ params: { id }}) => {
    const res = await serverClient.getGameDetails({ gameId: id })
    // TODO: handle res
    return (
        <GameWrapper>
            <GameBoard />
        </GameWrapper>
    )    
}

export default GamePage