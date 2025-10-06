import * as jspb from 'google-protobuf'

import * as game_v1_game_pb from '../../game/v1/game_pb'; // proto import: "game/v1/game.proto"


export class GetCurrentGameRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCurrentGameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCurrentGameRequest): GetCurrentGameRequest.AsObject;
  static serializeBinaryToWriter(message: GetCurrentGameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCurrentGameRequest;
  static deserializeBinaryFromReader(message: GetCurrentGameRequest, reader: jspb.BinaryReader): GetCurrentGameRequest;
}

export namespace GetCurrentGameRequest {
  export type AsObject = {
  }
}

export class GetCurrentGameResponse extends jspb.Message {
  getGame(): game_v1_game_pb.Game | undefined;
  setGame(value?: game_v1_game_pb.Game): GetCurrentGameResponse;
  hasGame(): boolean;
  clearGame(): GetCurrentGameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCurrentGameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCurrentGameResponse): GetCurrentGameResponse.AsObject;
  static serializeBinaryToWriter(message: GetCurrentGameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCurrentGameResponse;
  static deserializeBinaryFromReader(message: GetCurrentGameResponse, reader: jspb.BinaryReader): GetCurrentGameResponse;
}

export namespace GetCurrentGameResponse {
  export type AsObject = {
    game?: game_v1_game_pb.Game.AsObject,
  }
}

export class JoinGameRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinGameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: JoinGameRequest): JoinGameRequest.AsObject;
  static serializeBinaryToWriter(message: JoinGameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinGameRequest;
  static deserializeBinaryFromReader(message: JoinGameRequest, reader: jspb.BinaryReader): JoinGameRequest;
}

export namespace JoinGameRequest {
  export type AsObject = {
  }
}

export class JoinUpdate extends jspb.Message {
  getAcknowledgement(): JoinAcknowledgement | undefined;
  setAcknowledgement(value?: JoinAcknowledgement): JoinUpdate;
  hasAcknowledgement(): boolean;
  clearAcknowledgement(): JoinUpdate;

  getServerDetails(): ServerDetails | undefined;
  setServerDetails(value?: ServerDetails): JoinUpdate;
  hasServerDetails(): boolean;
  clearServerDetails(): JoinUpdate;

  getError(): JoinError | undefined;
  setError(value?: JoinError): JoinUpdate;
  hasError(): boolean;
  clearError(): JoinUpdate;

  getUpdateCase(): JoinUpdate.UpdateCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: JoinUpdate): JoinUpdate.AsObject;
  static serializeBinaryToWriter(message: JoinUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinUpdate;
  static deserializeBinaryFromReader(message: JoinUpdate, reader: jspb.BinaryReader): JoinUpdate;
}

export namespace JoinUpdate {
  export type AsObject = {
    acknowledgement?: JoinAcknowledgement.AsObject,
    serverDetails?: ServerDetails.AsObject,
    error?: JoinError.AsObject,
  }

  export enum UpdateCase { 
    UPDATE_NOT_SET = 0,
    ACKNOWLEDGEMENT = 1,
    SERVER_DETAILS = 2,
    ERROR = 3,
  }
}

export class JoinAcknowledgement extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): JoinAcknowledgement;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinAcknowledgement.AsObject;
  static toObject(includeInstance: boolean, msg: JoinAcknowledgement): JoinAcknowledgement.AsObject;
  static serializeBinaryToWriter(message: JoinAcknowledgement, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinAcknowledgement;
  static deserializeBinaryFromReader(message: JoinAcknowledgement, reader: jspb.BinaryReader): JoinAcknowledgement;
}

export namespace JoinAcknowledgement {
  export type AsObject = {
    message: string,
  }
}

export class ServerDetails extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): ServerDetails;

  getPort(): number;
  setPort(value: number): ServerDetails;

  getPlayerId(): string;
  setPlayerId(value: string): ServerDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ServerDetails): ServerDetails.AsObject;
  static serializeBinaryToWriter(message: ServerDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerDetails;
  static deserializeBinaryFromReader(message: ServerDetails, reader: jspb.BinaryReader): ServerDetails;
}

export namespace ServerDetails {
  export type AsObject = {
    address: string,
    port: number,
    playerId: string,
  }
}

export class JoinError extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): JoinError;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinError.AsObject;
  static toObject(includeInstance: boolean, msg: JoinError): JoinError.AsObject;
  static serializeBinaryToWriter(message: JoinError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinError;
  static deserializeBinaryFromReader(message: JoinError, reader: jspb.BinaryReader): JoinError;
}

export namespace JoinError {
  export type AsObject = {
    message: string,
  }
}

