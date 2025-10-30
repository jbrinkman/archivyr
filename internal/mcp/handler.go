// Package mcp implements the Model Context Protocol handlers for ruleset operations.
package mcp

import (
	"context"
	"fmt"

	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// Handler manages MCP protocol interactions for ruleset operations
type Handler struct {
	rulesetService ruleset.ServiceInterface
	server         *server.MCPServer
}

// NewHandler creates a new MCP handler with the given ruleset service
func NewHandler(service ruleset.ServiceInterface) *Handler {
	return &Handler{
		rulesetService: service,
	}
}

// Start initializes the MCP server with stdio transport and starts serving requests
func (h *Handler) Start() error {
	log.Info().Msg("Initializing MCP server")

	// Create MCP server with capabilities
	s := server.NewMCPServer(
		"MCP Ruleset Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	h.server = s

	log.Info().Msg("Registering resources")
	h.RegisterResources(s)

	log.Info().Msg("Registering tools")
	h.RegisterTools(s)

	log.Info().Msg("Starting MCP server with stdio transport")

	// Start server with stdio transport
	// This is a blocking call that handles MCP protocol communication
	if err := server.ServeStdio(s); err != nil {
		log.Error().Err(err).Msg("MCP server error")
		return fmt.Errorf("failed to serve stdio: %w", err)
	}

	log.Info().Msg("MCP server stopped")
	return nil
}

// RegisterResources registers ruleset resources with the MCP server
func (h *Handler) RegisterResources(s *server.MCPServer) {
	// Register resource template for ruleset retrieval by name
	resource := mcp.NewResource(
		"ruleset://{name}",
		"Ruleset",
		mcp.WithResourceDescription("AI editor ruleset with metadata and markdown content"),
		mcp.WithMIMEType("text/markdown"),
	)

	s.AddResource(resource, h.handleResourceRead)
}

// HandleResourceRead handles resource read requests for rulesets (exported for testing)
func (h *Handler) HandleResourceRead(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return h.handleResourceRead(ctx, req)
}

// handleResourceRead handles resource read requests for rulesets
func (h *Handler) handleResourceRead(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract ruleset name from URI
	// URI format: "ruleset://{name}" or "ruleset:{name}"
	uri := req.Params.URI
	name := extractNameFromURI(uri)

	if name == "" {
		return nil, fmt.Errorf("invalid URI format: %s", uri)
	}

	// Retrieve ruleset from service
	rs, err := h.rulesetService.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ruleset: %w", err)
	}

	// Format response with metadata and markdown content
	content := formatRulesetAsMarkdown(rs)

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "text/markdown",
			Text:     content,
		},
	}, nil
}

// extractNameFromURI extracts the ruleset name from the URI
// Supports formats: "ruleset://{name}" and "ruleset:{name}"
func extractNameFromURI(uri string) string {
	// Remove "ruleset://" prefix (11 characters)
	if len(uri) > 10 && uri[:10] == "ruleset://" {
		return uri[10:]
	}
	// Remove "ruleset:" prefix (8 characters)
	if len(uri) > 8 && uri[:8] == "ruleset:" {
		return uri[8:]
	}
	return ""
}

// formatRulesetAsMarkdown formats a ruleset with metadata as markdown
func formatRulesetAsMarkdown(rs *ruleset.Ruleset) string {
	// Format metadata header
	metadata := fmt.Sprintf(`---
name: %s
description: %s
tags: %v
created_at: %s
last_modified: %s
---

`, rs.Name, rs.Description, rs.Tags, rs.CreatedAt.Format("2006-01-02 15:04:05"), rs.LastModified.Format("2006-01-02 15:04:05"))

	// Append markdown content
	return metadata + rs.Markdown
}

// RegisterTools registers all CRUD tools with the MCP server
func (h *Handler) RegisterTools(s *server.MCPServer) {
	// Register upsert_ruleset tool (replaces create_ruleset and update_ruleset)
	upsertTool := mcp.NewTool("upsert_ruleset",
		mcp.WithDescription("Create a new ruleset or update an existing one. For new rulesets, all fields are required. For existing rulesets, only name is required and other fields are optional updates."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Snake_case ruleset name")),
		mcp.WithString("description", mcp.Description("Brief description of the ruleset (required for new rulesets)")),
		mcp.WithString("markdown", mcp.Description("Ruleset content in markdown format (required for new rulesets)")),
	)
	s.AddTool(upsertTool, h.handleUpsertRuleset)

	// Register get_ruleset tool
	getTool := mcp.NewTool("get_ruleset",
		mcp.WithDescription("Retrieve a ruleset by exact name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Exact ruleset name")),
	)
	s.AddTool(getTool, h.handleGetRuleset)

	// Register delete_ruleset tool
	deleteTool := mcp.NewTool("delete_ruleset",
		mcp.WithDescription("Delete a ruleset by name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Ruleset name to delete")),
	)
	s.AddTool(deleteTool, h.handleDeleteRuleset)

	// Register search_rulesets tool (replaces list_rulesets)
	searchTool := mcp.NewTool("search_rulesets",
		mcp.WithDescription("Search rulesets by name pattern. Omit pattern or use '*' to list all rulesets."),
		mcp.WithString("pattern", mcp.Description("Glob pattern (e.g., '*python*', 'style_*'). Defaults to '*' to list all rulesets.")),
	)
	s.AddTool(searchTool, h.handleSearchRulesets)
}

// HandleUpsertRuleset handles the upsert_ruleset tool invocation (exported for testing)
func (h *Handler) HandleUpsertRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleUpsertRuleset(ctx, req)
}

// handleUpsertRuleset handles the upsert_ruleset tool invocation
func (h *Handler) handleUpsertRuleset(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	// Extract optional parameters
	args := req.GetArguments()

	// Build ruleset struct for potential creation
	rs := &ruleset.Ruleset{
		Name: name,
	}

	// Build update struct for potential update
	updates := &ruleset.Update{}

	if description, ok := args["description"].(string); ok {
		rs.Description = description
		updates.Description = &description
	}

	if markdown, ok := args["markdown"].(string); ok {
		rs.Markdown = markdown
		updates.Markdown = &markdown
	}

	// Extract optional tags parameter
	if tagsParam, ok := args["tags"]; ok {
		if tagsList, ok := tagsParam.([]interface{}); ok {
			tags := make([]string, 0, len(tagsList))
			for _, tag := range tagsList {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
			rs.Tags = tags
			updates.Tags = &tags
		}
	} else {
		rs.Tags = []string{}
	}

	// Perform upsert
	err = h.rulesetService.Upsert(rs, updates)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to upsert ruleset: %v", err)), nil
	}

	// Check if it was a create or update to provide appropriate message
	exists, _ := h.rulesetService.Exists(name)
	if exists {
		return mcp.NewToolResultText(fmt.Sprintf("Successfully upserted ruleset '%s'", name)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully upserted ruleset '%s'", name)), nil
}

// HandleCreateRuleset handles the create_ruleset tool invocation (exported for testing).
//
// Deprecated: Use HandleUpsertRuleset instead.
func (h *Handler) HandleCreateRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleCreateRuleset(ctx, req)
}

// handleCreateRuleset handles the create_ruleset tool invocation.
//
// Deprecated: Use handleUpsertRuleset instead.
func (h *Handler) handleCreateRuleset(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameters
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	description, err := req.RequireString("description")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'description': %v", err)), nil
	}

	markdown, err := req.RequireString("markdown")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'markdown': %v", err)), nil
	}

	// Extract optional tags parameter
	tags := req.GetStringSlice("tags", []string{})

	// Create ruleset
	rs := &ruleset.Ruleset{
		Name:        name,
		Description: description,
		Tags:        tags,
		Markdown:    markdown,
	}

	err = h.rulesetService.Create(rs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create ruleset: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully created ruleset '%s'", name)), nil
}

// HandleGetRuleset handles the get_ruleset tool invocation (exported for testing)
func (h *Handler) HandleGetRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleGetRuleset(ctx, req)
}

// handleGetRuleset handles the get_ruleset tool invocation
func (h *Handler) handleGetRuleset(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	// Retrieve ruleset
	rs, err := h.rulesetService.Get(name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to retrieve ruleset: %v", err)), nil
	}

	// Format response
	content := formatRulesetAsMarkdown(rs)
	return mcp.NewToolResultText(content), nil
}

// HandleUpdateRuleset handles the update_ruleset tool invocation (exported for testing)
func (h *Handler) HandleUpdateRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleUpdateRuleset(ctx, req)
}

// handleUpdateRuleset handles the update_ruleset tool invocation
func (h *Handler) handleUpdateRuleset(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	// Build update struct with optional parameters
	updates := &ruleset.Update{}
	args := req.GetArguments()

	if description, ok := args["description"].(string); ok {
		updates.Description = &description
	}

	if markdown, ok := args["markdown"].(string); ok {
		updates.Markdown = &markdown
	}

	if tagsParam, ok := args["tags"]; ok {
		if tagsList, ok := tagsParam.([]interface{}); ok {
			tags := make([]string, 0, len(tagsList))
			for _, tag := range tagsList {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
			updates.Tags = &tags
		}
	}

	// Update ruleset
	err = h.rulesetService.Update(name, updates)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to update ruleset: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated ruleset '%s'", name)), nil
}

// HandleDeleteRuleset handles the delete_ruleset tool invocation (exported for testing)
func (h *Handler) HandleDeleteRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleDeleteRuleset(ctx, req)
}

// handleDeleteRuleset handles the delete_ruleset tool invocation
func (h *Handler) handleDeleteRuleset(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	// Delete ruleset
	err = h.rulesetService.Delete(name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete ruleset: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted ruleset '%s'", name)), nil
}

// HandleListRulesets handles the list_rulesets tool invocation (exported for testing)
func (h *Handler) HandleListRulesets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleListRulesets(ctx, req)
}

// handleListRulesets handles the list_rulesets tool invocation
func (h *Handler) handleListRulesets(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// List all rulesets
	rulesets, err := h.rulesetService.List()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list rulesets: %v", err)), nil
	}

	// Format response
	if len(rulesets) == 0 {
		return mcp.NewToolResultText("No rulesets found"), nil
	}

	result := fmt.Sprintf("Found %d ruleset(s):\n\n", len(rulesets))
	for _, rs := range rulesets {
		result += fmt.Sprintf("- **%s**: %s\n", rs.Name, rs.Description)
		if len(rs.Tags) > 0 {
			result += fmt.Sprintf("  Tags: %v\n", rs.Tags)
		}
		result += fmt.Sprintf("  Created: %s, Modified: %s\n\n",
			rs.CreatedAt.Format("2006-01-02 15:04:05"),
			rs.LastModified.Format("2006-01-02 15:04:05"))
	}

	return mcp.NewToolResultText(result), nil
}

// HandleSearchRulesets handles the search_rulesets tool invocation (exported for testing)
func (h *Handler) HandleSearchRulesets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.handleSearchRulesets(ctx, req)
}

// handleSearchRulesets handles the search_rulesets tool invocation
func (h *Handler) handleSearchRulesets(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract optional pattern parameter, default to "*" for listing all
	args := req.GetArguments()
	pattern := "*"
	if patternArg, ok := args["pattern"].(string); ok && patternArg != "" {
		pattern = patternArg
	}

	// Search rulesets
	rulesets, err := h.rulesetService.Search(pattern)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search rulesets: %v", err)), nil
	}

	// Format response
	if len(rulesets) == 0 {
		if pattern == "*" {
			return mcp.NewToolResultText("No rulesets found"), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("No rulesets found matching pattern '%s'", pattern)), nil
	}

	var result string
	if pattern == "*" {
		result = fmt.Sprintf("Found %d ruleset(s):\n\n", len(rulesets))
	} else {
		result = fmt.Sprintf("Found %d ruleset(s) matching '%s':\n\n", len(rulesets), pattern)
	}

	for _, rs := range rulesets {
		result += fmt.Sprintf("- **%s**: %s\n", rs.Name, rs.Description)
		if len(rs.Tags) > 0 {
			result += fmt.Sprintf("  Tags: %v\n", rs.Tags)
		}
		result += fmt.Sprintf("  Created: %s, Modified: %s\n\n",
			rs.CreatedAt.Format("2006-01-02 15:04:05"),
			rs.LastModified.Format("2006-01-02 15:04:05"))
	}

	return mcp.NewToolResultText(result), nil
}
