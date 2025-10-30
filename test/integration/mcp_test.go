package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jbrinkman/archivyr/internal/mcp"
	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/jbrinkman/archivyr/internal/valkey"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCPIntegration_ToolInvocations tests MCP tool invocations with the handler
func TestMCPIntegration_ToolInvocations(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client and service
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	service := ruleset.NewService(client)
	handler := mcp.NewHandler(service)

	ctx := context.Background()

	t.Run("CreateRuleset_Success", func(t *testing.T) {
		// Create tool request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "create_ruleset",
				Arguments: map[string]interface{}{
					"name":        "test_create_ruleset",
					"description": "Test ruleset for MCP integration",
					"tags":        []interface{}{"test", "mcp"},
					"markdown":    "# Test Ruleset\n\nThis is a test.",
				},
			},
		}

		// Invoke tool handler directly
		result, err := handler.HandleUpsertRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "Successfully upserted ruleset")
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "test_create_ruleset")

		// Verify ruleset was created
		rs, err := service.Get("test_create_ruleset")
		require.NoError(t, err)
		assert.Equal(t, "test_create_ruleset", rs.Name)
		assert.Equal(t, "Test ruleset for MCP integration", rs.Description)
	})

	t.Run("CreateRuleset_MissingParameter", func(t *testing.T) {
		// Create tool request with missing required parameter
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "create_ruleset",
				Arguments: map[string]interface{}{
					"description": "Missing name parameter",
					"markdown":    "# Test",
				},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleUpsertRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "missing required parameter 'name'")
	})

	t.Run("CreateRuleset_DuplicateName", func(t *testing.T) {
		// Create first ruleset
		rs := &ruleset.Ruleset{
			Name:        "duplicate_test",
			Description: "First ruleset",
			Tags:        []string{"test"},
			Markdown:    "# First",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Try to create duplicate
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "create_ruleset",
				Arguments: map[string]interface{}{
					"name":        "duplicate_test",
					"description": "Duplicate ruleset",
					"markdown":    "# Duplicate",
				},
			},
		}

		result, err := handler.HandleUpsertRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		// Upsert should succeed and update the existing ruleset
		assert.False(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "Successfully upserted ruleset")
	})

	t.Run("GetRuleset_Success", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "test_get_ruleset",
			Description: "Test get operation",
			Tags:        []string{"test", "get"},
			Markdown:    "# Get Test\n\nContent here.",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Create get request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "get_ruleset",
				Arguments: map[string]interface{}{
					"name": "test_get_ruleset",
				},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleGetRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		text := result.Content[0].(mcplib.TextContent).Text
		assert.Contains(t, text, "name: test_get_ruleset")
		assert.Contains(t, text, "description: Test get operation")
		assert.Contains(t, text, "# Get Test")
	})

	t.Run("GetRuleset_NotFound", func(t *testing.T) {
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "get_ruleset",
				Arguments: map[string]interface{}{
					"name": "nonexistent_ruleset",
				},
			},
		}

		result, err := handler.HandleGetRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "not found")
	})

	t.Run("UpdateRuleset_Success", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "test_update_ruleset",
			Description: "Original description",
			Tags:        []string{"test"},
			Markdown:    "# Original",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Wait to ensure timestamp difference
		time.Sleep(100 * time.Millisecond)

		// Create update request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "update_ruleset",
				Arguments: map[string]interface{}{
					"name":        "test_update_ruleset",
					"description": "Updated description",
					"tags":        []interface{}{"test", "updated"},
					"markdown":    "# Updated Content",
				},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleUpdateRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "Successfully updated")

		// Verify update
		updated, err := service.Get("test_update_ruleset")
		require.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
		assert.Equal(t, []string{"test", "updated"}, updated.Tags)
		assert.Equal(t, "# Updated Content", updated.Markdown)
	})

	t.Run("UpdateRuleset_PartialUpdate", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "test_partial_update",
			Description: "Original description",
			Tags:        []string{"test"},
			Markdown:    "# Original",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Update only description
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "update_ruleset",
				Arguments: map[string]interface{}{
					"name":        "test_partial_update",
					"description": "Only description updated",
				},
			},
		}

		result, err := handler.HandleUpdateRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		// Verify only description changed
		updated, err := service.Get("test_partial_update")
		require.NoError(t, err)
		assert.Equal(t, "Only description updated", updated.Description)
		assert.Equal(t, []string{"test"}, updated.Tags) // Tags unchanged
		assert.Equal(t, "# Original", updated.Markdown) // Markdown unchanged
	})

	t.Run("DeleteRuleset_Success", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "test_delete_ruleset",
			Description: "To be deleted",
			Tags:        []string{"test"},
			Markdown:    "# Delete Me",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Create delete request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "delete_ruleset",
				Arguments: map[string]interface{}{
					"name": "test_delete_ruleset",
				},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleDeleteRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "Successfully deleted")

		// Verify deletion
		_, err = service.Get("test_delete_ruleset")
		assert.Error(t, err)
	})

	t.Run("ListRulesets_Success", func(t *testing.T) {
		// Create multiple rulesets
		for i := 1; i <= 3; i++ {
			rs := &ruleset.Ruleset{
				Name:        "list_test_" + string(rune('0'+i)),
				Description: "List test ruleset",
				Tags:        []string{"test", "list"},
				Markdown:    "# List Test",
			}
			err := service.Create(rs)
			require.NoError(t, err)
		}

		// Create list request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name:      "list_rulesets",
				Arguments: map[string]interface{}{},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleListRulesets(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		text := result.Content[0].(mcplib.TextContent).Text
		assert.Contains(t, text, "Found")
		assert.Contains(t, text, "ruleset(s)")
	})

	t.Run("SearchRulesets_Success", func(t *testing.T) {
		// Create rulesets with specific pattern
		for i := 1; i <= 3; i++ {
			rs := &ruleset.Ruleset{
				Name:        "search_pattern_" + string(rune('0'+i)),
				Description: "Search test ruleset",
				Tags:        []string{"test", "search"},
				Markdown:    "# Search Test",
			}
			err := service.Create(rs)
			require.NoError(t, err)
		}

		// Create search request
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "search_rulesets",
				Arguments: map[string]interface{}{
					"pattern": "search_pattern_*",
				},
			},
		}

		// Invoke tool handler
		result, err := handler.HandleSearchRulesets(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		text := result.Content[0].(mcplib.TextContent).Text
		assert.Contains(t, text, "Found")
		assert.Contains(t, text, "matching")
		assert.Contains(t, text, "search_pattern_")
	})

	t.Run("SearchRulesets_NoMatches", func(t *testing.T) {
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "search_rulesets",
				Arguments: map[string]interface{}{
					"pattern": "nonexistent_pattern_*",
				},
			},
		}

		result, err := handler.HandleSearchRulesets(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "No rulesets found")
	})
}

// TestMCPIntegration_ResourceRetrieval tests MCP resource retrieval via URI
func TestMCPIntegration_ResourceRetrieval(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client and service
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	service := ruleset.NewService(client)
	handler := mcp.NewHandler(service)

	ctx := context.Background()

	t.Run("ResourceRead_Success", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "resource_test_ruleset",
			Description: "Test resource retrieval",
			Tags:        []string{"test", "resource"},
			Markdown:    "# Resource Test\n\nThis is resource content.",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Create resource read request with double slash URI
		req := mcplib.ReadResourceRequest{
			Params: mcplib.ReadResourceParams{
				URI: "ruleset://resource_test_ruleset",
			},
		}

		// Invoke resource handler
		contents, err := handler.HandleResourceRead(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, contents)
		assert.Len(t, contents, 1)

		textContent := contents[0].(mcplib.TextResourceContents)
		assert.Equal(t, "ruleset://resource_test_ruleset", textContent.URI)
		assert.Equal(t, "text/markdown", textContent.MIMEType)
		assert.Contains(t, textContent.Text, "name: resource_test_ruleset")
		assert.Contains(t, textContent.Text, "description: Test resource retrieval")
		assert.Contains(t, textContent.Text, "# Resource Test")
	})

	t.Run("ResourceRead_SingleColonURI", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "single_colon_test",
			Description: "Test single colon URI",
			Tags:        []string{"test"},
			Markdown:    "# Single Colon Test",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		// Create resource read request with single colon URI
		req := mcplib.ReadResourceRequest{
			Params: mcplib.ReadResourceParams{
				URI: "ruleset:single_colon_test",
			},
		}

		// Invoke resource handler
		contents, err := handler.HandleResourceRead(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, contents)
		assert.Len(t, contents, 1)

		textContent := contents[0].(mcplib.TextResourceContents)
		assert.Equal(t, "ruleset:single_colon_test", textContent.URI)
		assert.Contains(t, textContent.Text, "name: single_colon_test")
	})

	t.Run("ResourceRead_NotFound", func(t *testing.T) {
		req := mcplib.ReadResourceRequest{
			Params: mcplib.ReadResourceParams{
				URI: "ruleset://nonexistent_resource",
			},
		}

		// Invoke resource handler
		_, err := handler.HandleResourceRead(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ResourceRead_InvalidURI", func(t *testing.T) {
		req := mcplib.ReadResourceRequest{
			Params: mcplib.ReadResourceParams{
				URI: "invalid://uri/format",
			},
		}

		// Invoke resource handler
		_, err := handler.HandleResourceRead(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URI format")
	})

	t.Run("ResourceRead_EmptyURI", func(t *testing.T) {
		req := mcplib.ReadResourceRequest{
			Params: mcplib.ReadResourceParams{
				URI: "",
			},
		}

		// Invoke resource handler
		_, err := handler.HandleResourceRead(ctx, req)
		assert.Error(t, err)
	})
}

// TestMCPIntegration_ErrorResponses tests error response formatting
func TestMCPIntegration_ErrorResponses(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client and service
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	service := ruleset.NewService(client)
	handler := mcp.NewHandler(service)

	ctx := context.Background()

	t.Run("ErrorResponse_InvalidName", func(t *testing.T) {
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "create_ruleset",
				Arguments: map[string]interface{}{
					"name":        "Invalid-Name-With-Dashes",
					"description": "Test invalid name",
					"markdown":    "# Test",
				},
			},
		}

		result, err := handler.HandleUpsertRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)

		text := result.Content[0].(mcplib.TextContent).Text
		assert.Contains(t, text, "snake_case")
	})

	t.Run("ErrorResponse_UpdateNonexistent", func(t *testing.T) {
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "update_ruleset",
				Arguments: map[string]interface{}{
					"name":        "nonexistent_for_update",
					"description": "New description",
				},
			},
		}

		result, err := handler.HandleUpdateRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "not found")
	})

	t.Run("ErrorResponse_DeleteNonexistent", func(t *testing.T) {
		// Create a ruleset first so the error message includes existing rulesets
		rs := &ruleset.Ruleset{
			Name:        "error_test_ruleset",
			Description: "For error testing",
			Tags:        []string{"test"},
			Markdown:    "# Error Test",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "delete_ruleset",
				Arguments: map[string]interface{}{
					"name": "nonexistent_for_delete",
				},
			},
		}

		result, err := handler.HandleDeleteRuleset(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)

		text := result.Content[0].(mcplib.TextContent).Text
		assert.Contains(t, text, "not found")
		// Should include list of existing rulesets
		assert.Contains(t, text, "Existing rulesets")
	})

	t.Run("ErrorResponse_EmptySearchPattern", func(t *testing.T) {
		req := mcplib.CallToolRequest{
			Params: mcplib.CallToolParams{
				Name: "search_rulesets",
				Arguments: map[string]interface{}{
					"pattern": "",
				},
			},
		}

		result, err := handler.HandleSearchRulesets(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(mcplib.TextContent).Text, "pattern cannot be empty")
	})
}

// TestMCPIntegration_ConcurrentToolInvocations tests concurrent tool invocations
func TestMCPIntegration_ConcurrentToolInvocations(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client and service
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	service := ruleset.NewService(client)
	handler := mcp.NewHandler(service)

	ctx := context.Background()

	t.Run("ConcurrentCreates", func(t *testing.T) {
		numGoroutines := 10
		results := make(chan error, numGoroutines)

		for i := range numGoroutines {
			go func(index int) {
				req := mcplib.CallToolRequest{
					Params: mcplib.CallToolParams{
						Name: "create_ruleset",
						Arguments: map[string]interface{}{
							"name":        "concurrent_mcp_" + string(rune('0'+index)),
							"description": "Concurrent MCP test",
							"markdown":    "# Concurrent",
						},
					},
				}

				result, err := handler.HandleUpsertRuleset(ctx, req)
				if err != nil {
					results <- err
					return
				}
				if result.IsError {
					results <- assert.AnError
					return
				}
				results <- nil
			}(i)
		}

		// Collect results
		for range numGoroutines {
			err := <-results
			assert.NoError(t, err)
		}
	})

	t.Run("ConcurrentReads", func(t *testing.T) {
		// Create a ruleset first
		rs := &ruleset.Ruleset{
			Name:        "concurrent_read_test",
			Description: "Concurrent read test",
			Tags:        []string{"test"},
			Markdown:    "# Concurrent Read",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		numGoroutines := 20
		results := make(chan error, numGoroutines)

		for range numGoroutines {
			go func() {
				req := mcplib.CallToolRequest{
					Params: mcplib.CallToolParams{
						Name: "get_ruleset",
						Arguments: map[string]interface{}{
							"name": "concurrent_read_test",
						},
					},
				}

				result, err := handler.HandleGetRuleset(ctx, req)
				if err != nil {
					results <- err
					return
				}
				if result.IsError {
					results <- assert.AnError
					return
				}
				results <- nil
			}()
		}

		// Collect results
		for range numGoroutines {
			err := <-results
			assert.NoError(t, err)
		}
	})
}
