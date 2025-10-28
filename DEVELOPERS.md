# Developer Guide

## Prerequisites

- Go 1.21 or later
- Docker (for integration tests and deployment)
- golangci-lint (for code quality)

## Development Setup

```bash
# Clone the repository
git clone https://github.com/jbrinkman/archivyr.git
cd archivyr

# Download dependencies
go mod download

# Run tests
go test ./... -v -race -cover

# Run linter
golangci-lint run
```

## Project Structure

```
archivyr/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── valkey/            # Valkey client wrapper
│   ├── ruleset/           # Core business logic
│   ├── mcp/               # MCP protocol handlers
│   └── util/              # Shared utilities
├── test/                  # Test suites
│   ├── integration/       # Integration tests
│   └── e2e/               # End-to-end tests
└── docker/                # Docker configuration
```

## Building

```bash
# Build binary
go build -o bin/archivyr ./cmd/mcp-ruleset-server

# Build Docker image
docker build -f docker/Dockerfile -t archivyr:latest .
```

## Testing

```bash
# Unit tests
go test ./... -v

# With coverage
go test ./... -v -race -cover

# Integration tests
go test ./test/integration/... -v

# E2E tests
go test ./test/e2e/... -v
```

## Running Locally

```bash
# Set environment variables
export VALKEY_HOST=localhost
export VALKEY_PORT=6379
export LOG_LEVEL=debug

# Run the server
go run ./cmd/mcp-ruleset-server
```

## Architecture

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

## API Reference

See [API.md](docs/API.md) for MCP protocol details and tool/resource specifications.
