

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache git make

# Copiar módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binarios
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/consumer ./cmd/consumer

# Stage final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binarios
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/consumer .

# Copiar migraciones
# COPY --from=builder /app/migrations ./migrations

# # Exponer puerto
# EXPOSE 8080

# Default command (puede ser sobreescrito en docker-compose)
CMD ["./api"]