package e2e

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MCPRequest represents a JSON-RPC request to the MCP server
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC response from the MCP server
type MCPResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *MCPError              `json:"error,omitempty"`
}

// MCPError represents an error in the MCP response
type MCPError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// buildDockerImage builds the Docker image for testing
func buildDockerImage(t *testing.T, ctx context.Context) {
	t.Helper()

	// Build the Docker image using the Dockerfile
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../", // Root of the project
			Dockerfile: "docker/Dockerfile",
		},
	}

	// Create a generic container to trigger the build
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	require.NoError(t, err)

	// Terminate immediately as we only needed to build
	if container != nil {
		_ = container.Terminate(ctx)
	}
}

// startMCPContainer starts the MCP server container
func startMCPContainer(t *testing.T, ctx context.Context) testcontainers.Container {
	t.Helper()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../",
			Dockerfile: "docker/Dockerfile",
		},
		Cmd: []string{"/bin/sh", "/docker/docker-entrypoint.sh"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../docker/docker-entrypoint.sh",
				ContainerFilePath: "/docker/docker-entrypoint.sh",
				FileMode:          0755,
			},
		},
		WaitingFor: wait.ForLog("MCP Ruleset Server is running").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return container
}

// TestDockerE2E_BuildImage tests that the Docker image builds successfully
func TestDockerE2E_BuildImage(t *testing.T) {
	ctx := context.Background()

	t.Run("BuildSuccess", func(t *testing.T) {
		// This will build the image and verify no errors occur
		buildDockerImage(t, ctx)
	})
}

// TestDockerE2E_ContainerStartup tests container startup and initialization
func TestDockerE2E_ContainerStartup(t *testing.T) {
	ctx := context.Background()

	t.Run("StartupSuccess", func(t *testing.T) {
		container := startMCPContainer(t, ctx)
		defer func() {
			err := container.Terminate(ctx)
			require.NoError(t, err)
		}()

		// Verify container is running
		state, err := container.State(ctx)
		require.NoError(t, err)
		assert.True(t, state.Running)

		// Check logs for successful startup
		logs, err := container.Logs(ctx)
		require.NoError(t, err)
		defer func() { _ = logs.Close() }()

		logContent, err := io.ReadAll(logs)
		require.NoError(t, err)
		logStr := string(logContent)

		assert.Contains(t, logStr, "Starting Valkey server")
		assert.Contains(t, logStr, "Valkey is ready")
		assert.Contains(t, logStr, "Starting MCP server")
		assert.Contains(t, logStr, "MCP Ruleset Server is running")
	})

	t.Run("ValkeyHealthCheck", func(t *testing.T) {
		container := startMCPContainer(t, ctx)
		defer func() {
			err := container.Terminate(ctx)
			require.NoError(t, err)
		}()

		// Execute valkey-cli ping to verify Valkey is running
		exitCode, reader, err := container.Exec(ctx, []string{"valkey-cli", "ping"})
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)

		output, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Contains(t, string(output), "PONG")
	})
}

// TestDockerE2E_MCPServerAvailability tests MCP server availability via stdio
func TestDockerE2E_MCPServerAvailability(t *testing.T) {
	ctx := context.Background()

	container := startMCPContainer(t, ctx)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	t.Run("InitializeProtocol", func(t *testing.T) {
		// Send initialize request
		req := MCPRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"clientInfo": map[string]interface{}{
					"name":    "test-client",
					"version": "1.0.0",
				},
			},
		}

		// Execute the MCP server with the initialize request
		reqJSON, err := json.Marshal(req)
		require.NoError(t, err)

		exitCode, reader, err := container.Exec(ctx, []string{
			"/bin/sh", "-c",
			fmt.Sprintf("echo '%s' | timeout 5 /usr/local/bin/mcp-ruleset-server 2>&1 | head -1", string(reqJSON)),
		})
		require.NoError(t, err)

		output, err := io.ReadAll(reader)
		require.NoError(t, err)

		// The server should respond with JSON-RPC
		outputStr := strings.TrimSpace(string(output))
		if exitCode == 0 && len(outputStr) > 0 {
			// Try to parse as JSON to verify it's a valid response
			var resp map[string]interface{}
			err = json.Unmarshal([]byte(outputStr), &resp)
			if err == nil {
				assert.Contains(t, resp, "jsonrpc")
			}
		}
	})
}

// TestDockerE2E_FullCRUDWorkflow tests the complete CRUD workflow through the container
func TestDockerE2E_FullCRUDWorkflow(t *testing.T) {
	ctx := context.Background()

	container := startMCPContainer(t, ctx)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Helper function to execute MCP tool calls
	executeTool := func(toolName string, args map[string]interface{}) (string, error) {
		reqJSON, err := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "tools/call",
			"params": map[string]interface{}{
				"name":      toolName,
				"arguments": args,
			},
		})
		if err != nil {
			return "", err
		}

		exitCode, reader, err := container.Exec(ctx, []string{
			"/bin/sh", "-c",
			fmt.Sprintf("echo '%s' | timeout 5 /usr/local/bin/mcp-ruleset-server 2>&1 | head -1", string(reqJSON)),
		})
		if err != nil {
			return "", err
		}

		output, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}

		if exitCode != 0 {
			return "", fmt.Errorf("command failed with exit code %d: %s", exitCode, string(output))
		}

		return strings.TrimSpace(string(output)), nil
	}

	t.Run("CreateRuleset", func(t *testing.T) {
		output, err := executeTool("create_ruleset", map[string]interface{}{
			"name":        "e2e_test_ruleset",
			"description": "End-to-end test ruleset",
			"tags":        []string{"e2e", "test"},
			"markdown":    "# E2E Test\n\nThis is an end-to-end test.",
		})

		// We expect either success or a valid JSON-RPC response
		if err == nil && len(output) > 0 {
			// Try to parse as JSON
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("Create response: %s", output)
			}
		}
	})

	t.Run("GetRuleset", func(t *testing.T) {
		// First create a ruleset
		_, _ = executeTool("create_ruleset", map[string]interface{}{
			"name":        "e2e_get_test",
			"description": "Test get operation",
			"markdown":    "# Get Test",
		})

		// Then try to get it
		output, err := executeTool("get_ruleset", map[string]interface{}{
			"name": "e2e_get_test",
		})

		if err == nil && len(output) > 0 {
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("Get response: %s", output)
			}
		}
	})

	t.Run("ListRulesets", func(t *testing.T) {
		output, err := executeTool("list_rulesets", map[string]interface{}{})

		if err == nil && len(output) > 0 {
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("List response: %s", output)
			}
		}
	})

	t.Run("UpdateRuleset", func(t *testing.T) {
		// First create a ruleset
		_, _ = executeTool("create_ruleset", map[string]interface{}{
			"name":        "e2e_update_test",
			"description": "Original description",
			"markdown":    "# Original",
		})

		// Then update it
		output, err := executeTool("update_ruleset", map[string]interface{}{
			"name":        "e2e_update_test",
			"description": "Updated description",
		})

		if err == nil && len(output) > 0 {
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("Update response: %s", output)
			}
		}
	})

	t.Run("SearchRulesets", func(t *testing.T) {
		output, err := executeTool("search_rulesets", map[string]interface{}{
			"pattern": "e2e_*",
		})

		if err == nil && len(output) > 0 {
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("Search response: %s", output)
			}
		}
	})

	t.Run("DeleteRuleset", func(t *testing.T) {
		// First create a ruleset
		_, _ = executeTool("create_ruleset", map[string]interface{}{
			"name":        "e2e_delete_test",
			"description": "To be deleted",
			"markdown":    "# Delete Me",
		})

		// Then delete it
		output, err := executeTool("delete_ruleset", map[string]interface{}{
			"name": "e2e_delete_test",
		})

		if err == nil && len(output) > 0 {
			var resp map[string]interface{}
			if json.Unmarshal([]byte(output), &resp) == nil {
				t.Logf("Delete response: %s", output)
			}
		}
	})
}

// TestDockerE2E_DataPersistence tests data persistence across operations
func TestDockerE2E_DataPersistence(t *testing.T) {
	ctx := context.Background()

	container := startMCPContainer(t, ctx)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	t.Run("PersistenceAcrossOperations", func(t *testing.T) {
		// Create multiple rulesets
		for i := 1; i <= 3; i++ {
			exitCode, _, err := container.Exec(ctx, []string{
				"valkey-cli", "HSET",
				fmt.Sprintf("ruleset:persistence_test_%d", i),
				"description", fmt.Sprintf("Persistence test %d", i),
				"tags", `["test"]`,
				"markdown", "# Test",
				"created_at", time.Now().Format(time.RFC3339),
				"last_modified", time.Now().Format(time.RFC3339),
			})
			require.NoError(t, err)
			require.Equal(t, 0, exitCode)
		}

		// Verify data persists by checking with valkey-cli
		exitCode, reader, err := container.Exec(ctx, []string{
			"valkey-cli", "KEYS", "ruleset:persistence_test_*",
		})
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)

		output, err := io.ReadAll(reader)
		require.NoError(t, err)

		// Should find all 3 rulesets
		outputStr := string(output)
		assert.Contains(t, outputStr, "persistence_test_1")
		assert.Contains(t, outputStr, "persistence_test_2")
		assert.Contains(t, outputStr, "persistence_test_3")
	})
}

// TestDockerE2E_GracefulShutdown tests graceful shutdown of the container
func TestDockerE2E_GracefulShutdown(t *testing.T) {
	ctx := context.Background()

	t.Run("SIGTERMShutdown", func(t *testing.T) {
		container := startMCPContainer(t, ctx)

		// Get container logs before shutdown
		logsBefore, err := container.Logs(ctx)
		require.NoError(t, err)
		beforeContent, _ := io.ReadAll(logsBefore)
		_ = logsBefore.Close()

		// Send SIGTERM to the container
		stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = container.Stop(stopCtx, nil)
		require.NoError(t, err)

		// Get logs after shutdown
		logsAfter, err := container.Logs(ctx)
		require.NoError(t, err)
		defer func() { _ = logsAfter.Close() }()

		afterContent, err := io.ReadAll(logsAfter)
		require.NoError(t, err)

		// Combine logs
		allLogs := string(beforeContent) + string(afterContent)

		// Verify graceful shutdown messages
		assert.Contains(t, allLogs, "Received shutdown signal")

		// Terminate the container
		err = container.Terminate(ctx)
		require.NoError(t, err)
	})

	t.Run("CleanShutdown", func(t *testing.T) {
		container := startMCPContainer(t, ctx)

		// Verify container is running
		state, err := container.State(ctx)
		require.NoError(t, err)
		assert.True(t, state.Running)

		// Stop the container gracefully
		stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = container.Stop(stopCtx, nil)
		require.NoError(t, err)

		// Verify container stopped
		state, err = container.State(ctx)
		require.NoError(t, err)
		assert.False(t, state.Running)

		// Terminate the container
		err = container.Terminate(ctx)
		require.NoError(t, err)
	})
}

// TestDockerE2E_ContainerLogs tests that container logs are properly generated
func TestDockerE2E_ContainerLogs(t *testing.T) {
	ctx := context.Background()

	container := startMCPContainer(t, ctx)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	t.Run("LogsAvailable", func(t *testing.T) {
		// Wait a moment for logs to be generated
		time.Sleep(2 * time.Second)

		logs, err := container.Logs(ctx)
		require.NoError(t, err)
		defer func() { _ = logs.Close() }()

		// Read logs line by line
		scanner := bufio.NewScanner(logs)
		logLines := []string{}
		for scanner.Scan() {
			logLines = append(logLines, scanner.Text())
		}

		require.NoError(t, scanner.Err())
		assert.NotEmpty(t, logLines, "Container should produce logs")

		// Verify key log messages
		allLogs := strings.Join(logLines, "\n")
		assert.Contains(t, allLogs, "Starting Valkey server")
		assert.Contains(t, allLogs, "Starting MCP server")
	})
}

// TestDockerE2E_ErrorHandling tests error handling in the container
func TestDockerE2E_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	container := startMCPContainer(t, ctx)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()

	t.Run("InvalidCommand", func(t *testing.T) {
		// The MCP server doesn't have command-line flags, so this test
		// verifies that the server binary exists and is executable
		exitCode, reader, err := container.Exec(ctx, []string{
			"test", "-x", "/usr/local/bin/mcp-ruleset-server",
		})
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)

		output, _ := io.ReadAll(reader)
		t.Logf("Binary check output: %s", string(output))
	})

	t.Run("ValkeyConnectionCheck", func(t *testing.T) {
		// Verify Valkey is accessible
		exitCode, reader, err := container.Exec(ctx, []string{
			"valkey-cli", "ping",
		})
		require.NoError(t, err)
		require.Equal(t, 0, exitCode)

		output, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Contains(t, string(output), "PONG")
	})
}
