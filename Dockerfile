# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы модулей сначала для эффективного кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение (указываем правильный путь к main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/server .
# Копируем миграции
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]