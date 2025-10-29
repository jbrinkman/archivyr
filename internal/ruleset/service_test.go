package ruleset

import (
	"context"
	"testing"
	"time"

	"github.com/jbrinkman/archivyr/internal/valkey"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestValkey creates a Valkey container for testing
func setupTestValkey(t *testing.T) (*valkey.Client, func()) {
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

	port, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	client, err := valkey.NewClient(host, port.Port())
	require.NoError(t, err)

	cleanup := func() {
		_ = client.Close()
		_ = container.Terminate(ctx)
	}

	return client, cleanup
}

func TestCreate_Success(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	ruleset := &Ruleset{
		Name:        "test_ruleset",
		Description: "A test ruleset",
		Tags:        []string{"test", "example"},
		Markdown:    "# Test Ruleset\n\nThis is a test.",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	// Verify timestamps were set
	assert.False(t, ruleset.CreatedAt.IsZero())
	assert.False(t, ruleset.LastModified.IsZero())
	assert.Equal(t, ruleset.CreatedAt, ruleset.LastModified)

	// Verify the ruleset exists
	exists, err := service.Exists("test_ruleset")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestCreate_DuplicateName(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	ruleset := &Ruleset{
		Name:        "duplicate_test",
		Description: "First ruleset",
		Tags:        []string{"test"},
		Markdown:    "# First",
	}

	// Create first ruleset
	err := service.Create(ruleset)
	require.NoError(t, err)

	// Try to create duplicate
	duplicateRuleset := &Ruleset{
		Name:        "duplicate_test",
		Description: "Second ruleset",
		Tags:        []string{"test"},
		Markdown:    "# Second",
	}

	err = service.Create(duplicateRuleset)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Contains(t, err.Error(), "duplicate_test")
}

func TestCreate_InvalidName(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	testCases := []struct {
		name        string
		rulesetName string
	}{
		{"empty name", ""},
		{"uppercase", "TestRuleset"},
		{"spaces", "test ruleset"},
		{"hyphens", "test-ruleset"},
		{"starts with number", "1test"},
		{"starts with underscore", "_test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ruleset := &Ruleset{
				Name:        tc.rulesetName,
				Description: "Test",
				Tags:        []string{"test"},
				Markdown:    "# Test",
			}

			err := service.Create(ruleset)
			require.Error(t, err)
		})
	}
}

func TestCreate_TimestampSetting(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	beforeCreate := time.Now()

	ruleset := &Ruleset{
		Name:        "timestamp_test",
		Description: "Testing timestamps",
		Tags:        []string{"test"},
		Markdown:    "# Timestamp Test",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	afterCreate := time.Now()

	// Verify timestamps are within expected range
	assert.True(t, ruleset.CreatedAt.After(beforeCreate) || ruleset.CreatedAt.Equal(beforeCreate))
	assert.True(t, ruleset.CreatedAt.Before(afterCreate) || ruleset.CreatedAt.Equal(afterCreate))
	assert.Equal(t, ruleset.CreatedAt, ruleset.LastModified)
}

func TestExists(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Test non-existent ruleset
	exists, err := service.Exists("nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "exists_test",
		Description: "Test",
		Tags:        []string{"test"},
		Markdown:    "# Test",
	}

	err = service.Create(ruleset)
	require.NoError(t, err)

	// Test existing ruleset
	exists, err = service.Exists("exists_test")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestListNames(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Test empty list
	names, err := service.ListNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	// Create multiple rulesets
	rulesets := []string{"ruleset_one", "ruleset_two", "ruleset_three"}
	for _, name := range rulesets {
		ruleset := &Ruleset{
			Name:        name,
			Description: "Test",
			Tags:        []string{"test"},
			Markdown:    "# Test",
		}
		err := service.Create(ruleset)
		require.NoError(t, err)
	}

	// List all names
	names, err = service.ListNames()
	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.ElementsMatch(t, rulesets, names)
}
