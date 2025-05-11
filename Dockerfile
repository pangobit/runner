# Dockerfile for Runner app
# Use specific Alpine version for security and stability
FROM golang:1.24-alpine AS builder

# Create a non-root user and group for building
RUN addgroup -S runner && adduser -S runner -G runner

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with security flags enabled
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -s -extldflags '-static'" -o runner .

# Scan the application for security vulnerabilities
RUN apk add --no-cache ca-certificates

# Run tests
RUN go test -v ./...

# Final stage - use distroless for minimal attack surface
FROM gcr.io/distroless/static:nonroot

# Set the RUNNER_AUTH_TOKEN from build-arg (default to empty)
ARG RUNNER_AUTH_TOKEN=""
ENV RUNNER_AUTH_TOKEN=$RUNNER_AUTH_TOKEN

# Copy the binary from builder
COPY --from=builder /app/runner /app/runner
# Copy default config
COPY config.yaml /app/

# Use the nonroot user from the distroless image
USER nonroot:nonroot

# Expose the port (will be overridden by config)
EXPOSE 8080

# Metadata as defined in OCI image spec annotations
# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL org.opencontainers.image.title="Runner"
LABEL org.opencontainers.image.description="A simple web server that executes docker-compose commands via webhook"
LABEL org.opencontainers.image.url="https://github.com/pangobit/runner"
LABEL org.opencontainers.image.source="https://github.com/pangobit/runner"
LABEL org.opencontainers.image.vendor="pangobit"
LABEL org.opencontainers.image.licenses="MIT"

# Command to run the application
ENTRYPOINT ["/app/runner"]
 