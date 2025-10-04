import { MatchmakerClient, GetCurrentGameRequest } from '@globalfront/pb/matchmaker/v1/matchmaker';
import * as grpc from "@grpc/grpc-js";

const client = new MatchmakerClient('localhost:8080', grpc.credentials.createInsecure());

export async function getCurrentGame() {
    const req = new GetCurrentGameRequest()

    client.GetCurrentGame(req, (err, response) => {
        if (err) {
          console.error(err)
          return
        }
        console.log(response)
      })
}
