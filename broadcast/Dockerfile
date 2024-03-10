FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o broadcast ./main.go

CMD ["./broadcast"]