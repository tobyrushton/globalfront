import { matchmaker } from '@globalfront/pb/matchmaker/v1/matchmaker';
import { game } from '@globalfront/pb/game/v1/game';
import * as grpc from "@grpc/grpc-js";

const client = new matchmaker.v1.MatchmakerClient('localhost:4321', grpc.credentials.createInsecure());

export const getCurrentGame = async (): Promise<game.v1.Game>  => {
    const req = new matchmaker.v1.GetCurrentGameRequest()

    const gm: game.v1.Game = await new Promise(res => 
      client.GetCurrentGame(req, (err, response) => {
        if (err) {
          console.error(err)
          return
        } else {
          const gameObj = response?.toObject().game
          res(game.v1.Game.fromObject(gameObj!))
        }
      }))

    return gm
}
