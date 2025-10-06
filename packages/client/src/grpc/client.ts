import { GrpcTransport } from '@protobuf-ts/grpc-transport'
import { ChannelCredentials } from '@grpc/grpc-js'
import { MatchmakerClient } from '@globalfront/pb/matchmaker/v1/matchmaker.client'

const grpcTransport = new GrpcTransport({
    host: 'localhost:4321',
    channelCredentials: ChannelCredentials.createInsecure(),
}) 

export const serverClient = new MatchmakerClient(grpcTransport)
