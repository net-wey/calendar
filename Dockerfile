FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-calendar ./cmd/main

FROM alpine:latest
WORKDIR /app

# Копируем только необходимое
COPY --from=builder /app/task-calendar .
COPY --from=builder /app/local.yaml .
COPY --from=builder /app/internal/config ./internal/config

ENV CONFIG_PATH=/app/local.yaml  

CMD ["./task-calendar"]