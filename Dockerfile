# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate gqlgen code
RUN go run github.com/99designs/gqlgen generate

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ordersystem ./cmd/ordersystem

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy compiled binary
COPY --from=builder /app/ordersystem .

# Copy migrations (needed at runtime by golang-migrate file:// source)
COPY --from=builder /app/sql ./sql

EXPOSE 8000 50051 8080

CMD ["./ordersystem"]
