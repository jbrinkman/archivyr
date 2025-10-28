# Project Structure

## Directory Layout

```
mcp-ruleset-server/
├── cmd/                          # Application entry points
│   └── mcp-ruleset-server/       # Main server application
│       └── main.go               # Entry point with initialization
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   ├── config.go             # Config struct and loading
│   │   └── config_test.go        # Config tests
│   ├── valkey/                   # Valkey client wrapper
│   │   ├── client.go             # Client implementation
│   │   └── client_test.go        # Client tests
│   ├── ruleset/                  # Core business logic
│   │   ├── service.go            # CRUD operations
│   │   ├── service_test.go       # Service tests
│   │   └── types.go              # Data models
│   ├── mcp/                      # MCP protocol handlers
│   │   ├── handler.go            # Main handler
│   │   ├── handler_test.go       # Handler tests
│   │   ├── resources.go          # Resource registration
│   │   └── tools.go              # Tool registration
│   └── util/                     # Shared utilities
│       ├── validation.go         # Name validation
│       ├── validation_test.go    # Validation tests
│       └── patterns.go           # Pattern matching
├── test/                         # Test suites
│   ├── integration/              # Integration tests
│   │   ├── valkey_test.go        # Valkey integration
│   │   └── mcp_test.go           # MCP protocol integration
│   └── e2e/                      # End-to-end tests
│       └── docker_test.go        # Docker container tests
├── docker/                       # Docker configuration
│   ├── Dockerfile                # Multi-stage build
│   └── docker-entrypoint.sh      # Container startup script
├── .github/workflows/            # CI/CD pipelines
│   ├── ci.yml                    # Continuous integration
│   └── cd.yml                    # Continuous deployment
├── examples/                     # Usage examples
│   ├── mcp-config.json           # MCP client config
│   └── sample-rulesets/          # Example rulesets
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md           # Architecture details
│   └── API.md                    # API reference
└── Configuration files
    ├── go.mod                    # Go module definition
    ├── go.sum                    # Dependency checksums
    ├── .releaserc.json           # Semantic release config
    ├── .golangci.yml             # Linter configuration
    ├── README.md                 # Project overview
    ├── CONTRIBUTING.md           # Contribution guidelines
    └── CHANGELOG.md              # Version history
```

## Package Organization

### cmd/mcp-ruleset-server

Application entry point. Handles initialization, configuration loading, dependency wiring, and graceful shutdown.

### internal/config

Configuration management with environment variable loading and validation. Provides defaults for Valkey connection and logging.

### internal/valkey

Thin wrapper around valkey-glide client. Handles connection lifecycle and provides health check functionality.

### internal/ruleset

Core business logic for ruleset management. Implements CRUD operations, validation, and Valkey data mapping. Uses hash data structures with key pattern `ruleset:{name}`.

### internal/mcp

MCP protocol implementation. Registers resources (exact-match retrieval) and tools (CRUD operations). Handles stdio transport and protocol message routing.

### internal/util

Shared utilities for validation (snake_case enforcement), pattern matching (glob patterns), and timestamp handling.

## Naming Conventions

- **Packages**: Lowercase, single word (config, valkey, ruleset)
- **Files**: Lowercase with underscores (service_test.go)
- **Types**: PascalCase (Ruleset, Service, Client)
- **Functions**: PascalCase for exported, camelCase for private
- **Ruleset names**: snake_case enforced (python_style_guide)
- **Valkey keys**: Pattern `ruleset:{name}` (ruleset:python_style_guide)

## Testing Strategy

- Unit tests alongside source files (*_test.go)
- Integration tests in test/integration/ using testcontainers
- E2E tests in test/e2e/ for full Docker workflow
- Minimum 80% code coverage target
- All error paths must be tested
