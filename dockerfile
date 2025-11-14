# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both services with correct paths
RUN go build -o api-server ./cmd/api-server
RUN go build -o redirect-server ./cmd/redirect-server

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy both binaries from builder stage
COPY --from=builder /app/api-server .
COPY --from=builder /app/redirect-server .

# Copy web templates and static files
COPY --from=builder /app/web ./web

EXPOSE 8080 8081

# Default command (will be overridden by docker-compose)
CMD ["./api-server"]