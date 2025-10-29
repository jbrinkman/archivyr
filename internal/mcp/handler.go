package mcp

import (
	"context"
	"fmt"

	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handler manages MCP protocol interactions for ruleset operations
type Handler struct {
	rulesetService ruleset.ServiceInterface
}

// NewHandler creates a new MCP handler with the given ruleset service
func NewHandler(service ruleset.ServiceInterface) *Handler {
	return &Handler{
		rulesetService: service,
	}
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

// handleResourceRead handles resource read requests for rulesets
func (h *Handler) handleResourceRead(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
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
	// Register create_ruleset tool
	createTool := mcp.NewTool("create_ruleset",
		mcp.WithDescription("Create a new ruleset with metadata and content"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Snake_case ruleset name")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Brief description of the ruleset")),
		mcp.WithString("markdown", mcp.Required(), mcp.Description("Ruleset content in markdown format")),
	)
	s.AddTool(createTool, h.handleCreateRuleset)

	// Register get_ruleset tool
	getTool := mcp.NewTool("get_ruleset",
		mcp.WithDescription("Retrieve a ruleset by exact name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Exact ruleset name")),
	)
	s.AddTool(getTool, h.handleGetRuleset)

	// Register update_ruleset tool
	updateTool := mcp.NewTool("update_ruleset",
		mcp.WithDescription("Update an existing ruleset"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Ruleset name to update")),
		mcp.WithString("description", mcp.Description("New description (optional)")),
		mcp.WithString("markdown", mcp.Description("New content (optional)")),
	)
	s.AddTool(updateTool, h.handleUpdateRuleset)

	// Register delete_ruleset tool
	deleteTool := mcp.NewTool("delete_ruleset",
		mcp.WithDescription("Delete a ruleset by name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Ruleset name to delete")),
	)
	s.AddTool(deleteTool, h.handleDeleteRuleset)

	// Register list_rulesets tool
	listTool := mcp.NewTool("list_rulesets",
		mcp.WithDescription("List all available rulesets"),
	)
	s.AddTool(listTool, h.handleListRulesets)

	// Register search_rulesets tool
	searchTool := mcp.NewTool("search_rulesets",
		mcp.WithDescription("Search rulesets by name pattern"),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Glob pattern (e.g., '*python*', 'style_*')")),
	)
	s.AddTool(searchTool, h.handleSearchRulesets)
}

// handleCreateRuleset handles the create_ruleset tool invocation
func (h *Handler) handleCreateRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// handleGetRuleset handles the get_ruleset tool invocation
func (h *Handler) handleGetRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// handleUpdateRuleset handles the update_ruleset tool invocation
func (h *Handler) handleUpdateRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	name, err := req.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'name': %v", err)), nil
	}

	// Build update struct with optional parameters
	updates := &ruleset.RulesetUpdate{}
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

// handleDeleteRuleset handles the delete_ruleset tool invocation
func (h *Handler) handleDeleteRuleset(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// handleListRulesets handles the list_rulesets tool invocation
func (h *Handler) handleListRulesets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

// handleSearchRulesets handles the search_rulesets tool invocation
func (h *Handler) handleSearchRulesets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract required parameter
	pattern, err := req.RequireString("pattern")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'pattern': %v", err)), nil
	}

	// Search rulesets
	rulesets, err := h.rulesetService.Search(pattern)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to search rulesets: %v", err)), nil
	}

	// Format response
	if len(rulesets) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No rulesets found matching pattern '%s'", pattern)), nil
	}

	result := fmt.Sprintf("Found %d ruleset(s) matching '%s':\n\n", len(rulesets), pattern)
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
