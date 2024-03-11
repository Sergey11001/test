FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o broadcast ./cmd/main.go

ENV CONFIG_PATH="./config/broadcast.yaml"

CMD ["./broadcast"]