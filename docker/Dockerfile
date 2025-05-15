FROM golang:1.21 as builder

WORKDIR /app

# Копируем сначала только файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o task-calendar ./cmd/main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/task-calendar .
CMD ["./task-calendar"]