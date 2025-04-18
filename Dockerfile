# Stage 1: билдим бинарник
FROM golang:1.23.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# сначала копируем go.mod и go.sum (для кеширования зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# потом копируем остальной проект
COPY . .

# билдим основной бинарник (путь cmd/main.go, если нет директории api)
RUN go build -o notes-api ./cmd/main.go

# Stage 2: минимальный образ для запуска
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# копируем бинарник из первого этапа
COPY --from=builder /app/notes-api .

# указываем порт (для EXPOSE — не обязательно, но удобно)
EXPOSE 8080

CMD ["./notes-api"]
