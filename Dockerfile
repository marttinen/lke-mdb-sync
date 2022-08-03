FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
ADD . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/sync

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/sync .
CMD "/app/sync"
