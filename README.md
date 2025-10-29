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

This project uses [Task](https://taskfile.dev) for build automation. Install Task first:

```bash
# macOS
brew install go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# Or install via Go
go install github.com/go-task/task/v3/cmd/task@latest
```

### Common Development Tasks

```bash
# Set up development environment
task dev:setup

# Run tests (same as CI)
task test

# Run linter (same as CI)
task lint

# Build binary
task build

# Run full CI pipeline locally
task ci

# See all available tasks
task
```

### Quick Development Workflow

```bash
# Before committing
task verify                 # Format, lint, and test

# Run tests quickly during development
task test:quick

# Watch for changes and auto-test (requires entr)
task dev:watch
```

## Contributing

All commits must include a DCO signoff. Use `git commit -s` to automatically add the signoff.

## License

BSD-3-Clause - see [LICENSE](LICENSE) for details.
