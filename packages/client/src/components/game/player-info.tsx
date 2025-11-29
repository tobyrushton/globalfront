"use client"

import { useGame } from "./provider"
import { Slider } from "../ui/slider"

export const PlayerInfo = () => {
    const { player, attackPercentage, setAttackPercentage } = useGame()

    return (
        <div className="absolute flex flex-col gap-2 left-0 bottom-0 m-6 w-md text-foreground">
            <p>
                Troops: {player.troopCount}
            </p>
            <p>
                Attack Percentage: {attackPercentage}%
            </p>
            <Slider 
                defaultValue={[attackPercentage]}
                onValueChange={(value: number[]) => setAttackPercentage(value[0])}
                max={100}
                step={1}
            />
        </div>
    )
}