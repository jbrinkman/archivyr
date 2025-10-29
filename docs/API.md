# MCP Ruleset Server API Reference

## Overview

The MCP Ruleset Server implements the Model Context Protocol (MCP) to provide centralized storage and management of AI editor rulesets. This document describes all available resources, tools, request/response formats, and error handling.

## Protocol Information

- **Protocol**: Model Context Protocol (MCP)
- **Transport**: stdio (standard input/output)
- **Message Format**: JSON-RPC 2.0
- **Server Name**: MCP Ruleset Server
- **Server Version**: 1.0.0

## Capabilities

The server supports the following MCP capabilities:

- **Resources**: Read-only access to rulesets via URI
- **Tools**: CRUD operations for ruleset management
- **Logging**: Structured logging for debugging

## Resources

Resources provide read-only access to rulesets using exact-match URIs.

### Ruleset Resource

Retrieve a ruleset by its exact name using a URI.

**URI Template**: `ruleset://{name}`

**MIME Type**: `text/markdown`

**Description**: AI editor ruleset with metadata and markdown content

#### Request Format

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/read",
  "params": {
    "uri": "ruleset://python_style_guide"
  }
}
```

#### Response Format

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "contents": [
      {
        "uri": "ruleset://python_style_guide",
        "mimeType": "text/markdown",
        "text": "---\nname: python_style_guide\ndescription: Python coding style guidelines\ntags: [python, style, pep8]\ncreated_at: 2025-10-29 10:30:00\nlast_modified: 2025-10-29 10:30:00\n---\n\n# Python Style Guide\n\n## Naming Conventions\n..."
      }
    ]
  }
}
```

#### Error Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32603,
    "message": "failed to retrieve ruleset: ruleset 'nonexistent' not found"
  }
}
```

## Tools

Tools provide CRUD operations for managing rulesets.

### create_ruleset

Create a new ruleset with metadata and content.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Snake_case ruleset name (e.g., `python_style_guide`) |
| `description` | string | Yes | Brief description of the ruleset |
| `markdown` | string | Yes | Ruleset content in markdown format |
| `tags` | array of strings | No | Categorization tags (default: empty array) |

#### Request Example

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "create_ruleset",
    "arguments": {
      "name": "python_style_guide",
      "description": "Python coding style guidelines for team projects",
      "tags": ["python", "style", "pep8"],
      "markdown": "# Python Style Guide\n\n## Naming Conventions\n\n- Use snake_case for functions and variables\n- Use PascalCase for classes\n\n## Imports\n\n- Group imports: standard library, third-party, local\n- Use absolute imports when possible"
    }
  }
}
```

#### Success Response

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Successfully created ruleset 'python_style_guide'"
      }
    ]
  }
}
```

#### Error Response (Duplicate Name)

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "failed to create ruleset: ruleset 'python_style_guide' already exists. Please choose a different name. Existing rulesets: [python_style_guide, go_conventions, javascript_rules]"
      }
    ],
    "isError": true
  }
}
```

#### Error Response (Invalid Name)

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "failed to create ruleset: invalid ruleset name 'Python-Style': must use snake_case (lowercase letters, numbers, and underscores only)"
      }
    ],
    "isError": true
  }
}
```

---

### get_ruleset

Retrieve a ruleset by its exact name.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Exact ruleset name |

#### Request Example

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "get_ruleset",
    "arguments": {
      "name": "python_style_guide"
    }
  }
}
```

#### Success Response

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "---\nname: python_style_guide\ndescription: Python coding style guidelines for team projects\ntags: [python style pep8]\ncreated_at: 2025-10-29 10:30:00\nlast_modified: 2025-10-29 10:30:00\n---\n\n# Python Style Guide\n\n## Naming Conventions\n\n- Use snake_case for functions and variables\n- Use PascalCase for classes\n\n## Imports\n\n- Group imports: standard library, third-party, local\n- Use absolute imports when possible"
      }
    ]
  }
}
```

#### Error Response (Not Found)

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "failed to retrieve ruleset: ruleset 'nonexistent' not found"
      }
    ],
    "isError": true
  }
}
```

---

### update_ruleset

Update an existing ruleset. Only provided fields will be updated.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Ruleset name to update |
| `description` | string | No | New description |
| `tags` | array of strings | No | New tags (replaces existing tags) |
| `markdown` | string | No | New content |

**Note**: The `last_modified` timestamp is automatically updated. The `created_at` timestamp is preserved.

#### Request Example (Partial Update)

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "update_ruleset",
    "arguments": {
      "name": "python_style_guide",
      "description": "Updated Python coding style guidelines with PEP 8 compliance",
      "tags": ["python", "style", "pep8", "updated"]
    }
  }
}
```

#### Request Example (Full Update)

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "update_ruleset",
    "arguments": {
      "name": "python_style_guide",
      "description": "Comprehensive Python style guide",
      "tags": ["python", "style", "pep8", "formatting"],
      "markdown": "# Python Style Guide\n\n## Updated Content\n\nNew guidelines here..."
    }
  }
}
```

#### Success Response

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Successfully updated ruleset 'python_style_guide'"
      }
    ]
  }
}
```

#### Error Response (Not Found)

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "failed to update ruleset: ruleset 'nonexistent' not found"
      }
    ],
    "isError": true
  }
}
```

---

### delete_ruleset

Delete a ruleset by name.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Ruleset name to delete |

#### Request Example

```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "delete_ruleset",
    "arguments": {
      "name": "python_style_guide"
    }
  }
}
```

#### Success Response

```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Successfully deleted ruleset 'python_style_guide'"
      }
    ]
  }
}
```

#### Error Response (Not Found)

```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "failed to delete ruleset: ruleset 'nonexistent' not found. Existing rulesets: [go_conventions, javascript_rules, testing_best_practices]"
      }
    ],
    "isError": true
  }
}
```

---

### list_rulesets

List all available rulesets with metadata (excluding markdown content).

#### Parameters

None

#### Request Example

```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "list_rulesets",
    "arguments": {}
  }
}
```

#### Success Response (With Rulesets)

```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Found 3 ruleset(s):\n\n- **python_style_guide**: Python coding style guidelines for team projects\n  Tags: [python style pep8]\n  Created: 2025-10-29 10:30:00, Modified: 2025-10-29 15:45:00\n\n- **go_conventions**: Go coding conventions and best practices\n  Tags: [go conventions]\n  Created: 2025-10-29 11:00:00, Modified: 2025-10-29 11:00:00\n\n- **javascript_rules**: JavaScript and TypeScript style rules\n  Tags: [javascript typescript eslint]\n  Created: 2025-10-29 12:15:00, Modified: 2025-10-29 14:20:00\n\n"
      }
    ]
  }
}
```

#### Success Response (Empty)

```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "No rulesets found"
      }
    ]
  }
}
```

---

### search_rulesets

Search for rulesets by name pattern using glob syntax.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `pattern` | string | Yes | Glob pattern (e.g., `*python*`, `style_*`, `*_guide`) |

**Pattern Syntax**:

- `*` matches any sequence of characters
- `?` matches any single character
- Patterns are matched against the full ruleset name

#### Request Example

```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "search_rulesets",
    "arguments": {
      "pattern": "*python*"
    }
  }
}
```

#### Success Response (With Results)

```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Found 2 ruleset(s) matching '*python*':\n\n- **python_style_guide**: Python coding style guidelines for team projects\n  Tags: [python style pep8]\n  Created: 2025-10-29 10:30:00, Modified: 2025-10-29 15:45:00\n\n- **python_testing**: Python testing best practices\n  Tags: [python testing pytest]\n  Created: 2025-10-29 13:00:00, Modified: 2025-10-29 13:00:00\n\n"
      }
    ]
  }
}
```

#### Success Response (No Results)

```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "No rulesets found matching pattern 'nonexistent*'"
      }
    ]
  }
}
```

#### Additional Pattern Examples

```json
// Match all style guides
{"pattern": "*_style_*"}

// Match all guides
{"pattern": "*_guide"}

// Match specific prefix
{"pattern": "python_*"}

// Match all rulesets (equivalent to list_rulesets)
{"pattern": "*"}
```

---

## Error Handling

### Error Categories

The server returns errors in two formats depending on the error type:

1. **Protocol Errors**: Returned as JSON-RPC error objects
2. **Tool Errors**: Returned as tool results with `isError: true`

### Common Error Codes

| Code | Description | Example |
|------|-------------|---------|
| -32600 | Invalid Request | Malformed JSON-RPC request |
| -32601 | Method Not Found | Unknown MCP method |
| -32602 | Invalid Params | Missing required parameters |
| -32603 | Internal Error | Server-side errors (Valkey connection, etc.) |

### Validation Errors

#### Invalid Ruleset Name

Ruleset names must follow snake_case convention:

- Lowercase letters (a-z)
- Numbers (0-9)
- Underscores (_)
- Must not start or end with underscore
- Must not contain consecutive underscores

**Error Message Format**:

```
invalid ruleset name '{name}': must use snake_case (lowercase letters, numbers, and underscores only)
```

**Examples**:

- ✅ Valid: `python_style_guide`, `api_v2_rules`, `test123`
- ❌ Invalid: `Python-Style`, `api__rules`, `_private`, `style-guide`

#### Missing Required Parameters

**Error Message Format**:

```
missing required parameter '{parameter_name}': {details}
```

### Business Logic Errors

#### Duplicate Ruleset Name

When creating a ruleset with an existing name:

**Error Message Format**:

```
ruleset '{name}' already exists. Please choose a different name. Existing rulesets: [{list}]
```

#### Ruleset Not Found

When accessing a non-existent ruleset:

**Error Message Format**:

```
ruleset '{name}' not found
```

For delete operations, includes list of existing rulesets:

```
ruleset '{name}' not found. Existing rulesets: [{list}]
```

### Connection Errors

#### Valkey Connection Failure

When the server cannot connect to Valkey:

**Error Message Format**:

```
failed to connect to Valkey: {error_details}
```

**Note**: Connection errors typically occur during server startup and will cause the server to exit with a non-zero status code.

---

## Data Models

### Ruleset

Complete ruleset structure with all metadata and content.

```json
{
  "name": "python_style_guide",
  "description": "Python coding style guidelines for team projects",
  "tags": ["python", "style", "pep8"],
  "markdown": "# Python Style Guide\n\n## Content here...",
  "created_at": "2025-10-29T10:30:00Z",
  "last_modified": "2025-10-29T15:45:00Z"
}
```

**Fields**:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Snake_case identifier (unique) |
| `description` | string | Brief description of the ruleset |
| `tags` | array of strings | Categorization tags |
| `markdown` | string | Ruleset content in markdown format |
| `created_at` | timestamp | ISO 8601 timestamp of creation |
| `last_modified` | timestamp | ISO 8601 timestamp of last modification |

### Timestamp Format

All timestamps use RFC3339 format (ISO 8601):

- Format: `YYYY-MM-DDTHH:MM:SSZ`
- Timezone: UTC
- Example: `2025-10-29T10:30:00Z`

Display format in responses:

- Format: `YYYY-MM-DD HH:MM:SS`
- Example: `2025-10-29 10:30:00`

---

## Storage Schema

### Valkey Key Pattern

Rulesets are stored in Valkey using hash data structures:

**Key Pattern**: `ruleset:{name}`

**Example**: `ruleset:python_style_guide`

### Hash Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Ruleset description |
| `tags` | JSON string | Array of tags encoded as JSON |
| `markdown` | string | Markdown content |
| `created_at` | string | RFC3339 timestamp |
| `last_modified` | string | RFC3339 timestamp |

**Example Hash**:

```
Key: ruleset:python_style_guide
Fields:
  description: "Python coding style guidelines for team projects"
  tags: "[\"python\",\"style\",\"pep8\"]"
  markdown: "# Python Style Guide\n\n## Naming Conventions\n..."
  created_at: "2025-10-29T10:30:00Z"
  last_modified: "2025-10-29T15:45:00Z"
```

---

## Usage Examples

### Complete CRUD Workflow

```json
// 1. Create a ruleset
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_ruleset",
    "arguments": {
      "name": "api_design",
      "description": "RESTful API design principles",
      "tags": ["api", "rest", "design"],
      "markdown": "# API Design Principles\n\n## Resource Naming\n..."
    }
  }
}

// 2. List all rulesets
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "list_rulesets",
    "arguments": {}
  }
}

// 3. Get specific ruleset via resource
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/read",
  "params": {
    "uri": "ruleset://api_design"
  }
}

// 4. Update the ruleset
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "update_ruleset",
    "arguments": {
      "name": "api_design",
      "description": "RESTful API design principles and best practices",
      "tags": ["api", "rest", "design", "http"]
    }
  }
}

// 5. Search for API-related rulesets
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "search_rulesets",
    "arguments": {
      "pattern": "*api*"
    }
  }
}

// 6. Delete the ruleset
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "delete_ruleset",
    "arguments": {
      "name": "api_design"
    }
  }
}
```

### MCP Client Configuration

Example configuration for Claude Desktop or similar MCP clients:

```json
{
  "mcpServers": {
    "ruleset-server": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "mcp-ruleset-server:latest"
      ]
    }
  }
}
```

---

## Best Practices

### Naming Conventions

1. **Use descriptive names**: `python_style_guide` not `psg`
2. **Include context**: `api_design_principles` not just `design`
3. **Use consistent patterns**: `{language}_{topic}` or `{topic}_{type}`
4. **Avoid abbreviations**: Unless widely recognized (e.g., `api`, `http`)

### Tagging Strategy

1. **Use lowercase tags**: `python` not `Python`
2. **Be specific**: `["python", "pep8", "style"]` not just `["python"]`
3. **Include language/framework**: Always tag with primary technology
4. **Add category tags**: `style`, `testing`, `security`, `performance`

### Content Organization

1. **Use markdown headers**: Structure content with `#`, `##`, `###`
2. **Include examples**: Show code examples for guidelines
3. **Keep it concise**: Focus on actionable rules
4. **Update regularly**: Use update_ruleset to keep content current

### Error Handling

1. **Check for duplicates**: Before creating, consider searching first
2. **Validate names**: Ensure snake_case before calling create
3. **Handle not found**: Gracefully handle missing rulesets
4. **List on errors**: Use list_rulesets to discover available rulesets

---

## Limitations

### Current Version (1.0.0)

- **No authentication**: Server assumes trusted local environment
- **No authorization**: All clients have full CRUD access
- **No versioning**: Rulesets don't track version history
- **No bulk operations**: Must create/update/delete one at a time
- **No transactions**: Operations are not atomic across multiple rulesets
- **Pattern matching**: Basic glob patterns only (no regex)

### Scale Considerations

- **Recommended**: Up to 100 rulesets
- **Tested**: Couple dozen rulesets
- **Performance**: Optimized for small to medium collections

---

## Troubleshooting

### Common Issues

#### "Ruleset already exists"

**Cause**: Attempting to create a ruleset with a duplicate name

**Solution**:

1. Use `list_rulesets` to see existing names
2. Choose a different name
3. Or use `update_ruleset` to modify the existing ruleset

#### "Ruleset not found"

**Cause**: Accessing a non-existent ruleset

**Solution**:

1. Use `list_rulesets` to verify the name
2. Check for typos in the name
3. Ensure snake_case formatting

#### "Invalid ruleset name"

**Cause**: Name doesn't follow snake_case convention

**Solution**:

1. Use only lowercase letters, numbers, and underscores
2. Don't start or end with underscore
3. Avoid consecutive underscores
4. Examples: `my_ruleset`, `api_v2`, `test_123`

#### "Failed to connect to Valkey"

**Cause**: Valkey instance is not running or not accessible

**Solution**:

1. Verify Valkey is running: `valkey-cli ping`
2. Check connection settings: `VALKEY_HOST` and `VALKEY_PORT`
3. For Docker: Ensure container startup sequence is correct

---

## Version History

### 1.0.0 (Current)

- Initial release
- CRUD operations for rulesets
- Resource-based retrieval
- Pattern-based search
- Valkey-backed storage
- stdio transport
