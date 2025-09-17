FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY "cmd/voteflix" "cmd/voteflix"
COPY "cmd/voteflix-migrations" "cmd/voteflix-migrations"
COPY internal internal

RUN go build cmd/voteflix/main.go
RUN go build -o ./migrate cmd/voteflix-migrations/main.go

FROM golang:1.24 AS api

WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]

FROM golang:1.24 AS migrations

WORKDIR /app
COPY --from=build /app/migrate ./
COPY migrations migrations
CMD ["./migrate"]
