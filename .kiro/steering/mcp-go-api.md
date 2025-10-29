# MCP-Go API Reference

## Context7 Library ID

When you need to look up MCP-Go documentation, use this Context7-compatible library ID:

**Library ID:** `/mark3labs/mcp-go`

## Go Module

```go
import (
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)
```

## Quick API Reference

### Server Creation

```go
s := server.NewMCPServer(
    "Server Name",
    "1.0.0",
    server.WithToolCapabilities(true),
    server.WithResourceCapabilities(true, true),
    server.WithPromptCapabilities(true),
    server.WithLogging(),
)
```

### Adding Resources

```go
// Static resource
resource := mcp.NewResource(
    "scheme://{name}",  // URI template
    "Resource Name",
    mcp.WithResourceDescription("Description"),
    mcp.WithMIMEType("text/markdown"),
)

s.AddResource(resource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      req.Params.URI,
            MIMEType: "text/markdown",
            Text:     "content here",
        },
    }, nil
})
```

### Adding Tools

```go
tool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param1", mcp.Required(), mcp.Description("Parameter description")),
    mcp.WithNumber("param2", mcp.DefaultNumber(10)),
)

s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    param1, err := req.RequireString("param1")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    
    return mcp.NewToolResultText("result"), nil
})
```

### Starting Server

```go
// STDIO transport (for local use)
if err := server.ServeStdio(s); err != nil {
    log.Fatal(err)
}
```

## Important Notes

- Use `req.Params.URI` to access the resource URI in handlers
- Use `mcp.NewToolResultText()` for text responses
- Use `mcp.NewToolResultError()` for error responses
- Resource handlers return `[]mcp.ResourceContents`
- Tool handlers return `*mcp.CallToolResult`
