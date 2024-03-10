FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o service ./cmd/main.go

ENV CONFIG_PATH="./config/service1.yaml"

CMD ["./service"]