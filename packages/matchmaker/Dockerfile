FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY packages/matchmaker ./packages/matchmaker
COPY pb ./pb

WORKDIR /app/packages/matchmaker
RUN go build -o /matchmaker ./cmd

CMD ["/matchmaker"]
