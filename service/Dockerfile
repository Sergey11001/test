FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o service ./cmd/main.go

ENV ENV=local
ENV DIAL_TIMEOUT=4s
ENV BROADCAST_PORT=12345
ENV BROADCAST_PREFIX=aaa
ENV GRPC_TIMEOUT=1h
ENV CONFIG_PREFIX=S1

CMD ["./service"]