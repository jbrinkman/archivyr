# Archivyr

Centralized storage and management system for AI editor rulesets via the Model Context Protocol (MCP).

## Overview

Archivyr solves the problem of ruleset duplication and drift across multiple projects by providing a single source of truth for AI editor guidelines, rules, and steering documents. Store your rulesets once, access them everywhere.

## Features

- **Centralized Storage**: Store rulesets in a Valkey-backed system
- **MCP Protocol**: Access rulesets via the Model Context Protocol
- **CRUD Operations**: Create, read, update, and delete rulesets
- **Pattern Matching**: Search and list rulesets using glob patterns
- **Metadata Tracking**: Timestamps, tags, and descriptions for each ruleset
- **Docker Distribution**: Self-contained image with bundled Valkey instance

## Quick Start

```bash
# Run with Docker
docker run -i archivyr:latest

# Build from source
go build -o bin/archivyr ./cmd/mcp-ruleset-server
```

## Configuration

Configure via environment variables:

- `VALKEY_HOST`: Valkey host (default: localhost)
- `VALKEY_PORT`: Valkey port (default: 6379)
- `LOG_LEVEL`: Logging verbosity (default: info)

## Usage

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "archivyr": {
      "command": "docker",
      "args": ["run", "-i", "archivyr:latest"]
    }
  }
}
```

## Development

```bash
# Run tests
go test ./... -v -race -cover

# Run linter
golangci-lint run

# Build
go build -o bin/archivyr ./cmd/mcp-ruleset-server
```

## Contributing

All commits must include a DCO signoff. Use `git commit -s` to automatically add the signoff.

## License

BSD-3-Clause - see [LICENSE](LICENSE) for details.
