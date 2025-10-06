FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY packages/gamebox ./packages/gamebox
COPY pb ./pb

WORKDIR /app/packages/gamebox
RUN go build -o /gamebox ./cmd

CMD ["/gamebox"]