# Archivyr

Centralized storage and management system for AI editor rulesets via the Model Context Protocol (MCP).

## Overview

Archivyr solves the problem of ruleset duplication and drift across multiple projects by providing a single source of truth for AI editor guidelines, rules, and steering documents. Store your rulesets once, access them everywhere.

### The Problem

When working with AI editors across multiple projects, you often need to duplicate the same rulesets (coding standards, style guides, best practices) in each project. This leads to:

- Duplication of effort maintaining the same rules in multiple places
- Drift as rulesets evolve differently across projects
- Inconsistency in AI editor behavior across your projects

### The Solution

Archivyr provides a centralized MCP server that stores rulesets in Valkey and makes them available to any MCP-compatible AI editor. Create a ruleset once, reference it everywhere.

## Features

- **Centralized Storage**: Store rulesets in a Valkey-backed system
- **MCP Protocol**: Access rulesets via the Model Context Protocol
- **CRUD Operations**: Create, read, update, and delete rulesets
- **Pattern Matching**: Search and list rulesets using glob patterns
- **Metadata Tracking**: Timestamps, tags, and descriptions for each ruleset
- **Docker Distribution**: Self-contained image with bundled Valkey instance
- **Snake Case Naming**: Enforced naming convention for consistency

## Quick Start

### Using Docker (Recommended)

```bash
# Pull and run the latest image
docker run -i ghcr.io/joebrinkman/archivyr:latest

# Or build locally
docker build -f docker/Dockerfile -t archivyr:latest .
docker run -i archivyr:latest
```

### Building from Source

```bash
# Install dependencies
go mod download

# Build binary
task build

# Run (requires Valkey running on localhost:6379)
./bin/archivyr
```

## MCP Client Configuration

### Claude Desktop

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "archivyr": {
      "command": "docker",
      "args": ["run", "-i", "ghcr.io/joebrinkman/archivyr:latest"]
    }
  }
}
```

### Cursor

Add to your Cursor MCP configuration:

```json
{
  "mcpServers": {
    "archivyr": {
      "command": "docker",
      "args": ["run", "-i", "ghcr.io/joebrinkman/archivyr:latest"]
    }
  }
}
```

### Kiro

Add to your Kiro MCP configuration file:

**Workspace level**: `.kiro/settings/mcp.json`

**User level**: `~/.kiro/settings/mcp.json`

```json
{
  "mcpServers": {
    "archivyr": {
      "command": "docker",
      "args": ["run", "-i", "ghcr.io/joebrinkman/archivyr:latest"],
      "disabled": false,
      "autoApprove": []
    }
  }
}
```

### GitHub Copilot

Add to your GitHub Copilot MCP configuration:

**VS Code**: Settings → Extensions → GitHub Copilot → MCP Servers

```json
{
  "github.copilot.advanced": {
    "mcpServers": {
      "archivyr": {
        "command": "docker",
        "args": ["run", "-i", "ghcr.io/joebrinkman/archivyr:latest"]
      }
    }
  }
}
```

### Windsurf

Add to your Windsurf MCP configuration file:

**macOS**: `~/Library/Application Support/Windsurf/mcp_config.json`

**Windows**: `%APPDATA%\Windsurf\mcp_config.json`

**Linux**: `~/.config/Windsurf/mcp_config.json`

```json
{
  "mcpServers": {
    "archivyr": {
      "command": "docker",
      "args": ["run", "-i", "ghcr.io/joebrinkman/archivyr:latest"]
    }
  }
}
```

### Other MCP Clients

Archivyr uses stdio transport and works with any MCP-compatible client. See [examples/mcp-config.json](examples/mcp-config.json) for more configuration examples.

## Usage Examples

### Creating a Ruleset

Use the `create_ruleset` tool in your AI editor:

```
Create a new ruleset named "python_style_guide" with:
- Description: "Python coding standards for team projects"
- Tags: ["python", "style", "pep8"]
- Content: [your markdown content]
```

### Retrieving a Ruleset

Rulesets are exposed as MCP resources. Reference them by URI:

```
ruleset://python_style_guide
```

Or use the `get_ruleset` tool:

```
Get the ruleset named "python_style_guide"
```

### Searching Rulesets

Use the `search_rulesets` tool with glob patterns:

```
Search for rulesets matching "*python*"
Search for rulesets matching "style_*"
```

### Listing All Rulesets

Use the `list_rulesets` tool:

```
List all available rulesets
```

### Updating a Ruleset

Use the `update_ruleset` tool:

```
Update the ruleset "python_style_guide" with new description: "Updated Python standards"
```

### Deleting a Ruleset

Use the `delete_ruleset` tool:

```
Delete the ruleset named "old_ruleset"
```

## Available MCP Tools

- `create_ruleset`: Create a new ruleset with metadata and content
- `get_ruleset`: Retrieve a ruleset by exact name
- `update_ruleset`: Update an existing ruleset
- `delete_ruleset`: Delete a ruleset by name
- `list_rulesets`: List all available rulesets
- `search_rulesets`: Search rulesets by name pattern

## Available MCP Resources

- URI scheme: `ruleset://{name}`
- MIME type: `text/markdown`
- Example: `ruleset://python_style_guide`

## Configuration

Configure via environment variables:

- `VALKEY_HOST`: Valkey host (default: localhost)
- `VALKEY_PORT`: Valkey port (default: 6379)
- `LOG_LEVEL`: Logging verbosity (default: info)

## Architecture

Archivyr is built with:

- **Go 1.24+**: Core implementation
- **Valkey**: Key-value storage backend
- **valkey-glide v2**: Go client for Valkey
- **mcp-go**: MCP protocol implementation
- **Valkey Docker image**: Official Valkey image as base

### Data Storage

Rulesets are stored in Valkey as hashes with the key pattern `ruleset:{name}`:

```
Key: ruleset:python_style_guide
Fields:
  description: "Python coding standards"
  tags: ["python", "style", "pep8"]
  markdown: "# Python Style Guide\n..."
  created_at: "2025-10-28T10:30:00Z"
  last_modified: "2025-10-28T15:45:00Z"
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
