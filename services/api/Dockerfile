FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY "cmd/voteflix" "cmd/voteflix"
COPY internal internal

RUN go build cmd/voteflix/main.go

FROM golang:1.24-alpine

WORKDIR /app

COPY --from=build /app/main ./

CMD ["./main"]
