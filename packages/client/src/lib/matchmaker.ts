import { serverClient } from "@/grpc/client"
import { Game } from "@globalfront/pb/game/v1/game"

export const getCurrentGame = async (): Promise<Game>  => {
    const res = await serverClient.getCurrentGame({})

    return res.response.game!
  }
