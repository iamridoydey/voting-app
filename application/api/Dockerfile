# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

# Install git and other build tools if needed
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o server .

# Stage 2: Run the binary
FROM alpine:3.18 AS production

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose backend port
EXPOSE 3000

# Run the binary
CMD ["./server"]
