FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /shortener ./cmd/api

FROM alpine:latest as runner

WORKDIR /app

ENV REDIS_HOST=redis
ENV REDIS_PORT=6379

COPY --from=builder /shortener /app/shortener
EXPOSE 9000

CMD ["/app/shortener"]
