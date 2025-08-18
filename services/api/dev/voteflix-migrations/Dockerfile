FROM golang:1.24

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY "../../cmd" "cmd"
COPY ../../internal/app internal/app
COPY ../../migrations migrations

RUN go build -o voteflix-migrations cmd/voteflix-migrations/main.go

CMD ["./voteflix-migrations"]
