# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR internal/posts
COPY . .
WORKDIR posts
RUN go build -o main main.go

COPY posts/cfg.env .

ENV CGO_ENABLED=0

ENTRYPOINT ["./main"]