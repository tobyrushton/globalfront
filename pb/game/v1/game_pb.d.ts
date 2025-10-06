import * as jspb from 'google-protobuf'



export class Game extends jspb.Message {
  getId(): string;
  setId(value: string): Game;

  getPlayerCount(): number;
  setPlayerCount(value: number): Game;

  getMaxPlayers(): number;
  setMaxPlayers(value: number): Game;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Game.AsObject;
  static toObject(includeInstance: boolean, msg: Game): Game.AsObject;
  static serializeBinaryToWriter(message: Game, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Game;
  static deserializeBinaryFromReader(message: Game, reader: jspb.BinaryReader): Game;
}

export namespace Game {
  export type AsObject = {
    id: string,
    playerCount: number,
    maxPlayers: number,
  }
}

