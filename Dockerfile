# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o market-app .

# Final stage
FROM alpine:latest

# Устанавливаем только необходимое
RUN apk --no-cache add ca-certificates wget  # ← ДОБАВЬТЕ wget

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/market-app .

# Даем права на выполнение
RUN chmod +x market-app

EXPOSE 8070

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --retries=3 \
  CMD wget -q -O- http://localhost:8070/inventory/health || exit 1

CMD ["./market-app"]
