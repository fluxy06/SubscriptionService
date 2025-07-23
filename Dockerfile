# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o subscriptions_app main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscriptions_app .

CMD ["./subscriptions_app"]
