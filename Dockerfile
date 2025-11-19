# Multi-stage build for CLI Mahjong Engine
# Build stage optimized for compilation, runtime stage optimized for size

# ==============================================================================
# Stage 1: Build Stage
# ==============================================================================
FROM golang:1.21-alpine AS builder

# Install build dependencies
# - git: required for go mod operations
# - protobuf-dev: protoc compiler for protobuf schemas
# - make: build automation
# - ca-certificates: for HTTPS during go mod download
RUN apk add --no-cache git protobuf-dev make ca-certificates

WORKDIR /build

# Copy go mod files first for layer caching
# Docker caches layers, so if source code changes but go.mod doesn't,
# the go mod download step will be cached
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# - Builds the majhong_cli binary from cmd/majhong_cli
# - Output goes to bin/majhong_cli
RUN make build

# ==============================================================================
# Stage 2: Runtime Stage
# ==============================================================================
FROM alpine:latest

# Set headless mode for terminal applications
# TERM=dumb tells the application no terminal features are available
ENV TERM=dumb

# Add CA certificates for runtime HTTPS needs (if any)
RUN apk add --no-cache ca-certificates

# Create non-root user for security best practices
RUN addgroup -g 1000 majhong && adduser -D -u 1000 -G majhong majhong

# Copy binary from builder
COPY --from=builder /build/bin/majhong_cli /usr/local/bin/majhong_cli

# Verify binary was copied
RUN test -f /usr/local/bin/majhong_cli

# Set working directory
WORKDIR /data

# Switch to non-root user
USER majhong

# Default entrypoint and command
# Shows help by default when no arguments provided
ENTRYPOINT ["/usr/local/bin/majhong_cli"]
CMD ["--help"]
