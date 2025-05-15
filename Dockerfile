FROM golang:1.24

RUN apt-get update && apt-get install -y git curl

RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY . .

RUN go mod download


# Muda para o diretório onde está o main.go
WORKDIR /app/cmd/server


CMD ["air"]
