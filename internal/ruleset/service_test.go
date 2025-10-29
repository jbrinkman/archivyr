package ruleset

import (
	"context"
	"fmt"
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

func TestGet_Success(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset first
	ruleset := &Ruleset{
		Name:        "get_test",
		Description: "Test ruleset for Get",
		Tags:        []string{"test", "get"},
		Markdown:    "# Get Test\n\nThis is a test for Get operation.",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	// Retrieve the ruleset
	retrieved, err := service.Get("get_test")
	require.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Verify all fields
	assert.Equal(t, "get_test", retrieved.Name)
	assert.Equal(t, "Test ruleset for Get", retrieved.Description)
	assert.Equal(t, []string{"test", "get"}, retrieved.Tags)
	assert.Equal(t, "# Get Test\n\nThis is a test for Get operation.", retrieved.Markdown)
	assert.False(t, retrieved.CreatedAt.IsZero())
	assert.False(t, retrieved.LastModified.IsZero())
}

func TestGet_NotFound(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Try to get non-existent ruleset
	retrieved, err := service.Get("nonexistent_ruleset")
	require.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), "not found")
}

func TestGet_InvalidName(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Try to get with invalid name
	retrieved, err := service.Get("Invalid-Name")
	require.Error(t, err)
	assert.Nil(t, retrieved)
}

func TestList_Empty(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// List when no rulesets exist
	rulesets, err := service.List()
	require.NoError(t, err)
	assert.Empty(t, rulesets)
}

func TestList_WithRulesets(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create multiple rulesets
	testRulesets := []*Ruleset{
		{
			Name:        "list_test_one",
			Description: "First test ruleset",
			Tags:        []string{"test", "one"},
			Markdown:    "# First",
		},
		{
			Name:        "list_test_two",
			Description: "Second test ruleset",
			Tags:        []string{"test", "two"},
			Markdown:    "# Second",
		},
		{
			Name:        "list_test_three",
			Description: "Third test ruleset",
			Tags:        []string{"test", "three"},
			Markdown:    "# Third",
		},
	}

	for _, rs := range testRulesets {
		err := service.Create(rs)
		require.NoError(t, err)
	}

	// List all rulesets
	rulesets, err := service.List()
	require.NoError(t, err)
	assert.Len(t, rulesets, 3)

	// Verify all rulesets are present
	names := make([]string, len(rulesets))
	for i, rs := range rulesets {
		names[i] = rs.Name
	}
	assert.ElementsMatch(t, []string{"list_test_one", "list_test_two", "list_test_three"}, names)

	// Verify metadata is included
	for _, rs := range rulesets {
		assert.NotEmpty(t, rs.Description)
		assert.NotEmpty(t, rs.Tags)
		assert.NotEmpty(t, rs.Markdown)
		assert.False(t, rs.CreatedAt.IsZero())
		assert.False(t, rs.LastModified.IsZero())
	}
}

func TestSearch_WithWildcard(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create rulesets with different patterns
	testRulesets := []*Ruleset{
		{
			Name:        "python_style_guide",
			Description: "Python style guide",
			Tags:        []string{"python", "style"},
			Markdown:    "# Python Style",
		},
		{
			Name:        "python_testing_guide",
			Description: "Python testing guide",
			Tags:        []string{"python", "testing"},
			Markdown:    "# Python Testing",
		},
		{
			Name:        "javascript_style_guide",
			Description: "JavaScript style guide",
			Tags:        []string{"javascript", "style"},
			Markdown:    "# JavaScript Style",
		},
		{
			Name:        "go_conventions",
			Description: "Go conventions",
			Tags:        []string{"go", "conventions"},
			Markdown:    "# Go Conventions",
		},
	}

	for _, rs := range testRulesets {
		err := service.Create(rs)
		require.NoError(t, err)
	}

	// Search for python rulesets
	results, err := service.Search("python*")
	require.NoError(t, err)
	assert.Len(t, results, 2)

	names := make([]string, len(results))
	for i, rs := range results {
		names[i] = rs.Name
	}
	assert.ElementsMatch(t, []string{"python_style_guide", "python_testing_guide"}, names)
}

func TestSearch_WithSuffix(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create rulesets
	testRulesets := []*Ruleset{
		{
			Name:        "python_style_guide",
			Description: "Python style guide",
			Tags:        []string{"python"},
			Markdown:    "# Python",
		},
		{
			Name:        "javascript_style_guide",
			Description: "JavaScript style guide",
			Tags:        []string{"javascript"},
			Markdown:    "# JavaScript",
		},
		{
			Name:        "go_conventions",
			Description: "Go conventions",
			Tags:        []string{"go"},
			Markdown:    "# Go",
		},
	}

	for _, rs := range testRulesets {
		err := service.Create(rs)
		require.NoError(t, err)
	}

	// Search for style guides
	results, err := service.Search("*_style_guide")
	require.NoError(t, err)
	assert.Len(t, results, 2)

	names := make([]string, len(results))
	for i, rs := range results {
		names[i] = rs.Name
	}
	assert.ElementsMatch(t, []string{"python_style_guide", "javascript_style_guide"}, names)
}

func TestSearch_NoMatches(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "test_ruleset",
		Description: "Test",
		Tags:        []string{"test"},
		Markdown:    "# Test",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	// Search with pattern that doesn't match
	results, err := service.Search("nonexistent*")
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearch_EmptyPattern(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Search with empty pattern
	results, err := service.Search("")
	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "pattern cannot be empty")
}

func TestSearch_AllRulesets(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create multiple rulesets
	for i := 1; i <= 3; i++ {
		ruleset := &Ruleset{
			Name:        fmt.Sprintf("ruleset_%d", i),
			Description: fmt.Sprintf("Ruleset %d", i),
			Tags:        []string{"test"},
			Markdown:    fmt.Sprintf("# Ruleset %d", i),
		}
		err := service.Create(ruleset)
		require.NoError(t, err)
	}

	// Search with wildcard to get all
	results, err := service.Search("*")
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestUpdate_SuccessfulDescriptionUpdate(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "update_test",
		Description: "Original description",
		Tags:        []string{"test"},
		Markdown:    "# Original",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified
	time.Sleep(100 * time.Millisecond) // Ensure timestamp difference

	// Update description
	newDescription := "Updated description"
	updates := &RulesetUpdate{
		Description: &newDescription,
	}

	err = service.Update("update_test", updates)
	require.NoError(t, err)

	// Verify update
	updated, err := service.Get("update_test")
	require.NoError(t, err)
	assert.Equal(t, "Updated description", updated.Description)
	assert.Equal(t, []string{"test"}, updated.Tags) // Unchanged
	assert.Equal(t, "# Original", updated.Markdown) // Unchanged
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())
	// last_modified should be >= original (may be same second due to RFC3339 precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
}

func TestUpdate_SuccessfulTagsUpdate(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "tags_update_test",
		Description: "Test description",
		Tags:        []string{"old", "tags"},
		Markdown:    "# Test",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified
	time.Sleep(100 * time.Millisecond)

	// Update tags
	newTags := []string{"new", "updated", "tags"}
	updates := &RulesetUpdate{
		Tags: &newTags,
	}

	err = service.Update("tags_update_test", updates)
	require.NoError(t, err)

	// Verify update
	updated, err := service.Get("tags_update_test")
	require.NoError(t, err)
	assert.Equal(t, "Test description", updated.Description) // Unchanged
	assert.Equal(t, []string{"new", "updated", "tags"}, updated.Tags)
	assert.Equal(t, "# Test", updated.Markdown) // Unchanged
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())
	// last_modified should be >= original (may be same second due to RFC3339 precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
}

func TestUpdate_SuccessfulMarkdownUpdate(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "markdown_update_test",
		Description: "Test description",
		Tags:        []string{"test"},
		Markdown:    "# Original Content",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified
	time.Sleep(100 * time.Millisecond)

	// Update markdown
	newMarkdown := "# Updated Content\n\nThis is the new content."
	updates := &RulesetUpdate{
		Markdown: &newMarkdown,
	}

	err = service.Update("markdown_update_test", updates)
	require.NoError(t, err)

	// Verify update
	updated, err := service.Get("markdown_update_test")
	require.NoError(t, err)
	assert.Equal(t, "Test description", updated.Description) // Unchanged
	assert.Equal(t, []string{"test"}, updated.Tags)          // Unchanged
	assert.Equal(t, "# Updated Content\n\nThis is the new content.", updated.Markdown)
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())
	// last_modified should be >= original (may be same second due to RFC3339 precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
}

func TestUpdate_PartialUpdate(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "partial_update_test",
		Description: "Original description",
		Tags:        []string{"original", "tags"},
		Markdown:    "# Original",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified
	time.Sleep(100 * time.Millisecond)

	// Update only description and markdown, leave tags unchanged
	newDescription := "Updated description"
	newMarkdown := "# Updated"
	updates := &RulesetUpdate{
		Description: &newDescription,
		Markdown:    &newMarkdown,
	}

	err = service.Update("partial_update_test", updates)
	require.NoError(t, err)

	// Verify update
	updated, err := service.Get("partial_update_test")
	require.NoError(t, err)
	assert.Equal(t, "Updated description", updated.Description)
	assert.Equal(t, []string{"original", "tags"}, updated.Tags) // Unchanged
	assert.Equal(t, "# Updated", updated.Markdown)
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())
	// last_modified should be >= original (may be same second due to RFC3339 precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
}

func TestUpdate_AllFields(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "all_fields_update_test",
		Description: "Original description",
		Tags:        []string{"original"},
		Markdown:    "# Original",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified
	time.Sleep(100 * time.Millisecond)

	// Update all fields
	newDescription := "Completely updated description"
	newTags := []string{"updated", "all", "fields"}
	newMarkdown := "# Completely Updated\n\nAll fields changed."
	updates := &RulesetUpdate{
		Description: &newDescription,
		Tags:        &newTags,
		Markdown:    &newMarkdown,
	}

	err = service.Update("all_fields_update_test", updates)
	require.NoError(t, err)

	// Verify update
	updated, err := service.Get("all_fields_update_test")
	require.NoError(t, err)
	assert.Equal(t, "Completely updated description", updated.Description)
	assert.Equal(t, []string{"updated", "all", "fields"}, updated.Tags)
	assert.Equal(t, "# Completely Updated\n\nAll fields changed.", updated.Markdown)
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())
	// last_modified should be >= original (may be same second due to RFC3339 precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
}

func TestUpdate_TimestampHandling(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "timestamp_update_test",
		Description: "Test",
		Tags:        []string{"test"},
		Markdown:    "# Test",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	originalCreatedAt := ruleset.CreatedAt
	originalLastModified := ruleset.LastModified

	// Wait to ensure timestamp difference
	time.Sleep(100 * time.Millisecond)

	// Update
	newDescription := "Updated"
	updates := &RulesetUpdate{
		Description: &newDescription,
	}

	err = service.Update("timestamp_update_test", updates)
	require.NoError(t, err)

	// Verify timestamps
	updated, err := service.Get("timestamp_update_test")
	require.NoError(t, err)

	// created_at should be preserved
	assert.Equal(t, originalCreatedAt.Unix(), updated.CreatedAt.Unix())

	// last_modified should be updated (>= due to RFC3339 second precision)
	assert.True(t, updated.LastModified.Unix() >= originalLastModified.Unix())
	assert.True(t, updated.LastModified.Unix() >= originalCreatedAt.Unix())
}

func TestUpdate_NonExistentRuleset(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Try to update non-existent ruleset
	newDescription := "Updated description"
	updates := &RulesetUpdate{
		Description: &newDescription,
	}

	err := service.Update("nonexistent_ruleset", updates)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdate_InvalidName(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Try to update with invalid name
	newDescription := "Updated description"
	updates := &RulesetUpdate{
		Description: &newDescription,
	}

	err := service.Update("Invalid-Name", updates)
	require.Error(t, err)
}

func TestUpdate_EmptyUpdate(t *testing.T) {
	client, cleanup := setupTestValkey(t)
	defer cleanup()

	service := NewService(client)

	// Create a ruleset
	ruleset := &Ruleset{
		Name:        "empty_update_test",
		Description: "Original",
		Tags:        []string{"test"},
		Markdown:    "# Original",
	}

	err := service.Create(ruleset)
	require.NoError(t, err)

	// Update with no fields (should succeed but not change anything)
	updates := &RulesetUpdate{}

	err = service.Update("empty_update_test", updates)
	require.NoError(t, err)

	// Verify nothing changed except potentially last_modified
	updated, err := service.Get("empty_update_test")
	require.NoError(t, err)
	assert.Equal(t, "Original", updated.Description)
	assert.Equal(t, []string{"test"}, updated.Tags)
	assert.Equal(t, "# Original", updated.Markdown)
}
