# Technology Stack

## Language & Runtime

- **Go 1.21+**: Primary implementation language
- **Alpine Linux**: Docker base image for minimal footprint

## Core Dependencies

- **valkey-glide v2**: Go client library for Valkey interactions
- **mcp-go** (mark3labs/mcp-go): MCP protocol implementation
- **zerolog**: Structured logging
- **testify**: Testing utilities
- **testcontainers-go**: Integration testing with containers

## Data Store

- **Valkey**: Key-value store for ruleset persistence
- Data structure: Hashes with key pattern `ruleset:{name}`
- Default connection: localhost:6379

## Communication

- **stdio transport**: Standard input/output for MCP protocol
- **JSON-RPC**: MCP protocol message format

## Build & Development

### Common Commands

```bash
# Initialize module
go mod download

# Run tests
go test ./... -v -race -cover

# Run linter
golangci-lint run

# Build binary
go build -o bin/mcp-ruleset-server ./cmd/mcp-ruleset-server

# Build Docker image
docker build -f docker/Dockerfile -t mcp-ruleset-server:latest .

# Run container
docker run -i mcp-ruleset-server:latest

# Run integration tests
go test ./test/integration/... -v

# Run e2e tests
go test ./test/e2e/... -v
```

## CI/CD

- **GitHub Actions**: Automated workflows
- **golangci-lint**: Code quality enforcement
- **semantic-release**: Automated versioning and releases
- **Conventional Commits**: Commit message format enforcement

## Configuration

Environment variables:

- `VALKEY_HOST`: Valkey host (default: localhost)
- `VALKEY_PORT`: Valkey port (default: 6379)
- `LOG_LEVEL`: Logging verbosity (default: info)
