FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod download

RUN go build -o server .

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config.env ./config.env

RUN chmod +x ./server

EXPOSE 8080

ENTRYPOINT ["./server"]
