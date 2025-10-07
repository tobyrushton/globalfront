import { FC } from "react";
import { GameWrapper } from "@/components/game/wrapper";
import { GameBoard } from "@/components/game/board";

const GamePage: FC = () => {
    // do something to fetch the url here
    return (
        <GameWrapper>
            <GameBoard />
        </GameWrapper>
    )    
}

export default GamePage