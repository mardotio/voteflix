FROM golang:1.24

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download
COPY ../main.go .
COPY ../src src

ENTRYPOINT exec CompileDaemon -command="./api"