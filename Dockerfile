FROM golang:1.24.0 AS builder

WORKDIR /app

COPY . .

RUN go build -o chopper ./cmd/main.go

CMD ["./chopper"]