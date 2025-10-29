# Stage 1: Build Go binary
FROM golang:1.24-bookworm AS builder

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled (required for valkey-glide)
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=1 go build \
    -ldflags="-w -s" \
    -o mcp-ruleset-server \
    ./cmd/mcp-ruleset-server

# Stage 2: Runtime image based on official Valkey image
FROM valkey/valkey:9

# Set working directory
WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /build/mcp-ruleset-server /usr/local/bin/mcp-ruleset-server

# Make binary executable
RUN chmod +x /usr/local/bin/mcp-ruleset-server

# Set environment variables with defaults
ENV VALKEY_HOST=localhost \
    VALKEY_PORT=6379 \
    LOG_LEVEL=info

# Expose Valkey port (for potential future use)
EXPOSE 6379

# Volume for Valkey data persistence
VOLUME ["/data"]

# The entrypoint script will be added in the next task (task 17)
# For now, we'll use a simple command that starts the MCP server
# (Valkey startup will be handled by the entrypoint script in task 17)
CMD ["/usr/local/bin/mcp-ruleset-server"]
