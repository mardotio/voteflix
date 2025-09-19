FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY "cmd/voteflix" "cmd/voteflix"
COPY "cmd/voteflix-migrations" "cmd/voteflix-migrations"
COPY internal internal

RUN go build cmd/voteflix/main.go
RUN go build -o ./migrate cmd/voteflix-migrations/main.go

FROM golang:1.24-alpine AS api

WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]

FROM golang:1.24-alpine AS migrations

WORKDIR /app
COPY --from=build /app/migrate ./
COPY migrations migrations
CMD ["./migrate"]
