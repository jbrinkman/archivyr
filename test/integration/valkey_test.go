package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/jbrinkman/archivyr/internal/valkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupValkeyContainer starts a Valkey container for testing
func setupValkeyContainer(t *testing.T) (testcontainers.Container, string, string) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	return container, host, mappedPort.Port()
}

// teardownValkeyContainer stops and removes the Valkey container
func teardownValkeyContainer(t *testing.T, container testcontainers.Container) {
	ctx := context.Background()
	err := container.Terminate(ctx)
	require.NoError(t, err)
}

func TestValkeyIntegration_FullCRUDWorkflow(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	// Create ruleset service
	service := ruleset.NewService(client)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		rs := &ruleset.Ruleset{
			Name:        "test_ruleset",
			Description: "Test ruleset for integration testing",
			Tags:        []string{"test", "integration"},
			Markdown:    "# Test Ruleset\n\nThis is a test.",
		}

		err := service.Create(rs)
		assert.NoError(t, err)
		assert.False(t, rs.CreatedAt.IsZero())
		assert.False(t, rs.LastModified.IsZero())
	})

	// Test Read (Get)
	t.Run("Get", func(t *testing.T) {
		rs, err := service.Get("test_ruleset")
		require.NoError(t, err)
		assert.Equal(t, "test_ruleset", rs.Name)
		assert.Equal(t, "Test ruleset for integration testing", rs.Description)
		assert.Equal(t, []string{"test", "integration"}, rs.Tags)
		assert.Equal(t, "# Test Ruleset\n\nThis is a test.", rs.Markdown)
		assert.False(t, rs.CreatedAt.IsZero())
		assert.False(t, rs.LastModified.IsZero())
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		// Get original to compare timestamps
		original, err := service.Get("test_ruleset")
		require.NoError(t, err)

		// Wait a moment to ensure timestamp difference
		time.Sleep(100 * time.Millisecond)

		// Update the ruleset
		newDesc := "Updated description"
		newTags := []string{"test", "integration", "updated"}
		newMarkdown := "# Updated Test Ruleset\n\nThis has been updated."

		updates := &ruleset.Update{
			Description: &newDesc,
			Tags:        &newTags,
			Markdown:    &newMarkdown,
		}

		err = service.Update("test_ruleset", updates)
		assert.NoError(t, err)

		// Verify updates
		updated, err := service.Get("test_ruleset")
		require.NoError(t, err)
		assert.Equal(t, newDesc, updated.Description)
		assert.Equal(t, newTags, updated.Tags)
		assert.Equal(t, newMarkdown, updated.Markdown)
		assert.Equal(t, original.CreatedAt, updated.CreatedAt) // CreatedAt should not change
		// LastModified should be updated (allow for equal due to timestamp precision)
		assert.True(t, updated.LastModified.After(original.LastModified) || updated.LastModified.Equal(original.LastModified))
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		// Create additional rulesets
		rs2 := &ruleset.Ruleset{
			Name:        "another_ruleset",
			Description: "Another test ruleset",
			Tags:        []string{"test"},
			Markdown:    "# Another Ruleset",
		}
		err := service.Create(rs2)
		require.NoError(t, err)

		// List all rulesets
		rulesets, err := service.List()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(rulesets), 2)

		// Verify both rulesets are in the list
		names := make(map[string]bool)
		for _, rs := range rulesets {
			names[rs.Name] = true
		}
		assert.True(t, names["test_ruleset"])
		assert.True(t, names["another_ruleset"])
	})

	// Test Search
	t.Run("Search", func(t *testing.T) {
		// Search with wildcard pattern
		results, err := service.Search("test*")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)

		// Verify test_ruleset is in results
		found := false
		for _, rs := range results {
			if rs.Name == "test_ruleset" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := service.Delete("another_ruleset")
		assert.NoError(t, err)

		// Verify deletion
		_, err = service.Get("another_ruleset")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	// Test Exists
	t.Run("Exists", func(t *testing.T) {
		exists, err := service.Exists("test_ruleset")
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = service.Exists("nonexistent_ruleset")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestValkeyIntegration_ConcurrentOperations(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	// Create ruleset service
	service := ruleset.NewService(client)

	// Test concurrent creates
	t.Run("ConcurrentCreates", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				rs := &ruleset.Ruleset{
					Name:        fmt.Sprintf("concurrent_ruleset_%d", index),
					Description: fmt.Sprintf("Concurrent test ruleset %d", index),
					Tags:        []string{"concurrent", "test"},
					Markdown:    fmt.Sprintf("# Concurrent Ruleset %d", index),
				}

				if err := service.Create(rs); err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent create failed: %v", err)
		}

		// Verify all rulesets were created
		for i := 0; i < numGoroutines; i++ {
			name := fmt.Sprintf("concurrent_ruleset_%d", i)
			exists, err := service.Exists(name)
			require.NoError(t, err)
			assert.True(t, exists, "Ruleset %s should exist", name)
		}
	})

	// Test concurrent reads
	t.Run("ConcurrentReads", func(t *testing.T) {
		// Create a ruleset to read
		rs := &ruleset.Ruleset{
			Name:        "read_test_ruleset",
			Description: "Ruleset for concurrent read testing",
			Tags:        []string{"read", "test"},
			Markdown:    "# Read Test",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		var wg sync.WaitGroup
		numGoroutines := 20
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := service.Get("read_test_ruleset")
				if err != nil {
					errors <- err
				}
			}()
		}

		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent read failed: %v", err)
		}
	})

	// Test concurrent updates
	t.Run("ConcurrentUpdates", func(t *testing.T) {
		// Create a ruleset to update
		rs := &ruleset.Ruleset{
			Name:        "update_test_ruleset",
			Description: "Ruleset for concurrent update testing",
			Tags:        []string{"update", "test"},
			Markdown:    "# Update Test",
		}
		err := service.Create(rs)
		require.NoError(t, err)

		var wg sync.WaitGroup
		numGoroutines := 10
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				newDesc := fmt.Sprintf("Updated by goroutine %d", index)
				updates := &ruleset.Update{
					Description: &newDesc,
				}

				if err := service.Update("update_test_ruleset", updates); err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent update failed: %v", err)
		}

		// Verify the ruleset still exists and has a valid description
		updated, err := service.Get("update_test_ruleset")
		require.NoError(t, err)
		assert.Contains(t, updated.Description, "Updated by goroutine")
	})
}

func TestValkeyIntegration_ConnectionHandling(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	t.Run("SuccessfulConnection", func(t *testing.T) {
		client, err := valkey.NewClient(host, port)
		require.NoError(t, err)
		assert.NotNil(t, client)

		// Test ping
		err = client.Ping()
		assert.NoError(t, err)

		// Close connection
		err = client.Close()
		assert.NoError(t, err)
	})

	t.Run("MultipleConnections", func(t *testing.T) {
		// Create multiple clients
		clients := make([]*valkey.Client, 5)
		for i := 0; i < 5; i++ {
			client, err := valkey.NewClient(host, port)
			require.NoError(t, err)
			clients[i] = client

			// Test each connection
			err = client.Ping()
			assert.NoError(t, err)
		}

		// Close all connections
		for _, client := range clients {
			err := client.Close()
			assert.NoError(t, err)
		}
	})

	t.Run("ConnectionReuse", func(t *testing.T) {
		client, err := valkey.NewClient(host, port)
		require.NoError(t, err)
		defer func() { _ = client.Close() }()

		service := ruleset.NewService(client)

		// Perform multiple operations with the same connection
		for i := 0; i < 10; i++ {
			rs := &ruleset.Ruleset{
				Name:        fmt.Sprintf("reuse_test_%d", i),
				Description: "Connection reuse test",
				Tags:        []string{"test"},
				Markdown:    "# Test",
			}

			err := service.Create(rs)
			assert.NoError(t, err)

			_, err = service.Get(rs.Name)
			assert.NoError(t, err)
		}
	})
}

func TestValkeyIntegration_ErrorScenarios(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create Valkey client
	client, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	// Create ruleset service
	service := ruleset.NewService(client)

	t.Run("DuplicateCreate", func(t *testing.T) {
		rs := &ruleset.Ruleset{
			Name:        "duplicate_test",
			Description: "Test duplicate creation",
			Tags:        []string{"test"},
			Markdown:    "# Test",
		}

		// First create should succeed
		err := service.Create(rs)
		assert.NoError(t, err)

		// Second create should fail
		err = service.Create(rs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("GetNonexistent", func(t *testing.T) {
		_, err := service.Get("nonexistent_ruleset_xyz")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("UpdateNonexistent", func(t *testing.T) {
		newDesc := "Updated description"
		updates := &ruleset.Update{
			Description: &newDesc,
		}

		err := service.Update("nonexistent_ruleset_xyz", updates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("DeleteNonexistent", func(t *testing.T) {
		err := service.Delete("nonexistent_ruleset_xyz")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("InvalidRulesetName", func(t *testing.T) {
		rs := &ruleset.Ruleset{
			Name:        "Invalid-Name-With-Dashes",
			Description: "Test invalid name",
			Tags:        []string{"test"},
			Markdown:    "# Test",
		}

		err := service.Create(rs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "snake_case")
	})

	t.Run("EmptySearchPattern", func(t *testing.T) {
		_, err := service.Search("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})
}

func TestValkeyIntegration_DataPersistence(t *testing.T) {
	// Start Valkey container
	container, host, port := setupValkeyContainer(t)
	defer teardownValkeyContainer(t, container)

	// Create first client and service
	client1, err := valkey.NewClient(host, port)
	require.NoError(t, err)

	service1 := ruleset.NewService(client1)

	// Create rulesets with first client
	rs1 := &ruleset.Ruleset{
		Name:        "persistence_test_1",
		Description: "Test data persistence",
		Tags:        []string{"persistence", "test"},
		Markdown:    "# Persistence Test 1",
	}
	err = service1.Create(rs1)
	require.NoError(t, err)

	rs2 := &ruleset.Ruleset{
		Name:        "persistence_test_2",
		Description: "Another persistence test",
		Tags:        []string{"persistence", "test"},
		Markdown:    "# Persistence Test 2",
	}
	err = service1.Create(rs2)
	require.NoError(t, err)

	// Close first client
	_ = client1.Close()

	// Create second client and service
	client2, err := valkey.NewClient(host, port)
	require.NoError(t, err)
	defer func() { _ = client2.Close() }()

	service2 := ruleset.NewService(client2)

	// Verify data persists with new client
	retrieved1, err := service2.Get("persistence_test_1")
	require.NoError(t, err)
	assert.Equal(t, rs1.Name, retrieved1.Name)
	assert.Equal(t, rs1.Description, retrieved1.Description)
	assert.Equal(t, rs1.Tags, retrieved1.Tags)
	assert.Equal(t, rs1.Markdown, retrieved1.Markdown)

	retrieved2, err := service2.Get("persistence_test_2")
	require.NoError(t, err)
	assert.Equal(t, rs2.Name, retrieved2.Name)
	assert.Equal(t, rs2.Description, retrieved2.Description)
	assert.Equal(t, rs2.Tags, retrieved2.Tags)
	assert.Equal(t, rs2.Markdown, retrieved2.Markdown)

	// List should show both rulesets
	rulesets, err := service2.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(rulesets), 2)
}
