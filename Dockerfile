

# Base image with Go and necessary tools
FROM golang:1.25-alpine AS base
WORKDIR /app
RUN apk add --no-cache git make


# Build the Go binaries
FROM base AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/consumer ./cmd/consumer

# Produce a minimal image
FROM alpine:3.22
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/consumer .
CMD ["./api"]