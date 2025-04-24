# Dockerfile for Runner app
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o runner .

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/runner .

# Install Docker CLI and docker-compose
RUN apk add --no-cache docker-cli docker-compose

# Expose the port (will be overridden by config)
EXPOSE 8080

# Define volumes for config and compose files
VOLUME /app/config
VOLUME /app/compose

# Note: This container includes Docker CLI and docker-compose to interact with the host's Docker daemon
# via the mounted Docker socket (/var/run/docker.sock).

# Command to run the application
# Users can override the config path with -config flag if needed
CMD ["./runner", "-config", "/app/config/config.yaml"] 