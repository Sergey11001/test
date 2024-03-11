FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o broadcast ./cmd/main.go

ENV BROADCAST_IP="255.255.255.255"
ENV BROADCAST_PORT=12345
ENV BROADCAST_PREFIX="aaa"

CMD ["./broadcast"]