import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport'
import { MatchmakerClient } from '@globalfront/pb/matchmaker/v1/matchmaker.client'

const webTransport = new GrpcWebFetchTransport({
    baseUrl: 'http://localhost:8080',
    format: "binary",
})

export const webClient = new MatchmakerClient(webTransport)