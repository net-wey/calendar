FROM golang:1.21 as builder

WORKDIR /app

# 1. Сначала копируем только файлы модулей для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# 2. Копируем весь остальной код
COPY . .

# 3. Собираем приложение, указывая путь к main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o task-calendar ./cmd/main

# 4. Финальный легковесный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем только необходимые файлы
COPY --from=builder /app/task-calendar .
COPY --from=builder /app/init/init.sql ./init/  

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["./task-calendar"]