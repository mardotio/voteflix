FROM golang:1.24

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY "../../cmd" "cmd"
COPY ../../internal internal

ENTRYPOINT exec CompileDaemon -build="go build -o voteflix cmd/voteflix/main.go" -command="./voteflix"