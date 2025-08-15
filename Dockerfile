# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for SSL connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create certs directory with proper permissions
RUN mkdir -p /certs && chmod 755 /certs

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Expose both HTTP and HTTPS ports
EXPOSE 8080 8443

CMD ["./server"]
