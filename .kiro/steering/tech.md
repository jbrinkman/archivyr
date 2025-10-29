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

### Task Runner

This project uses [Task](https://taskfile.dev) to standardize build, test, and lint commands. **Always use Task commands instead of direct Go commands** to ensure consistency with CI and avoid skipping important steps.

### Common Task Commands

```bash
# Show all available tasks
task

# Development setup
task dev:setup              # Set up development environment
task deps                   # Download dependencies

# Code quality
task fmt                    # Format code
task lint                   # Run linter (same as CI)
task lint:fix               # Run linter with auto-fix

# Testing
task test                   # Run all tests with race detector and coverage (same as CI)
task test:unit              # Run unit tests only
task test:integration       # Run integration tests
task test:e2e               # Run end-to-end tests
task test:coverage          # Generate HTML coverage report
task test:quick             # Run tests without race detector (faster)

# Building
task build                  # Build binary
task build:all              # Build for multiple platforms

# Docker
task docker:build           # Build Docker image
task docker:run             # Run Docker container
task docker:test            # Run Docker smoke tests (same as CI)

# CI simulation
task ci                     # Run full CI pipeline locally
task ci:quick               # Run quick CI checks (lint + unit tests)

# Verification
task verify                 # Verify code is ready for commit (format, lint, test)

# Cleanup
task clean                  # Clean build artifacts
task clean:all              # Clean everything including Docker images
```

### Direct Go Commands (Avoid)

While direct Go commands work, they may skip important flags or steps used in CI. Use Task commands instead:

```bash
# ❌ Don't use these directly
go test ./...
go build ./cmd/mcp-ruleset-server
golangci-lint run

# ✅ Use Task commands instead
task test
task build
task lint
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
