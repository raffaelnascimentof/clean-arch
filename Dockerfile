# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git protobuf protobuf-dev

WORKDIR /app

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Regenerate protobuf files
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/infra/grpc/pb/order.proto

# Install gqlgen and generate GraphQL code
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
