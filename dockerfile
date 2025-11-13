# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder stage
COPY --from=builder /app/main .
# Copy web templates and static files
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./main"]