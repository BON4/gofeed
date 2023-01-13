# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR internal/accounts
COPY . .
WORKDIR accounts
RUN go build -o main main.go

COPY accounts/cfg.env .

ENV CGO_ENABLED=0

ENTRYPOINT ["./main" ]