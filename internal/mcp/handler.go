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
	rulesetService *ruleset.Service
	server         *server.MCPServer
}

// NewHandler creates a new MCP handler with the given ruleset service
func NewHandler(service *ruleset.Service) *Handler {
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
	// Remove "ruleset://" or "ruleset:" prefix
	if len(uri) > 11 && uri[:11] == "ruleset://" {
		return uri[11:]
	}
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
