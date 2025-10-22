# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gateway .

# Copy example keys file (users should mount their own)
COPY keys.example.json ./keys.example.json

# Expose port 8080
EXPOSE 8080

# Set environment variables with defaults
ENV SERVER_PORT=8080
ENV KEYS_FILE_PATH=/app/keys.json
ENV LOG_TO_FILE=false
ENV QUOTA_ENABLED=true
ENV QUOTA_LIMIT=100
ENV REQUEST_TIMEOUT=30

# Run the application
CMD ["./gateway"]
