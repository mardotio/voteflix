FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY "cmd/voteflix" "cmd/voteflix"
COPY internal internal

RUN go build cmd/voteflix/main.go
RUN go build -o ./migrate cmd/voteflix-migrations/main.go

FROM golang:1.24 as api

WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]

FROM golang:1.24 as migrations

WORKDIR /app
COPY --from=build /app/migrate ./
COPY migrations migrations
CMD ["./migrate"]
