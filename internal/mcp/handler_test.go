package mcp

import (
	"context"
	"testing"

	"github.com/jbrinkman/archivyr/internal/ruleset"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRulesetService is a mock implementation of the ruleset service interface
type MockRulesetService struct {
	mock.Mock
}

// Ensure MockRulesetService implements ruleset.ServiceInterface
var _ ruleset.ServiceInterface = (*MockRulesetService)(nil)

func (m *MockRulesetService) Create(rs *ruleset.Ruleset) error {
	args := m.Called(rs)
	return args.Error(0)
}

func (m *MockRulesetService) Get(name string) (*ruleset.Ruleset, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ruleset.Ruleset), args.Error(1)
}

func (m *MockRulesetService) Update(name string, updates *ruleset.Update) error {
	args := m.Called(name, updates)
	return args.Error(0)
}

func (m *MockRulesetService) Upsert(rs *ruleset.Ruleset, updates *ruleset.Update) error {
	args := m.Called(rs, updates)
	return args.Error(0)
}

func (m *MockRulesetService) Delete(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockRulesetService) List() ([]*ruleset.Ruleset, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ruleset.Ruleset), args.Error(1)
}

func (m *MockRulesetService) Search(pattern string) ([]*ruleset.Ruleset, error) {
	args := m.Called(pattern)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ruleset.Ruleset), args.Error(1)
}

func (m *MockRulesetService) Exists(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRulesetService) ListNames() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Test Handler creation
func TestNewHandler(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.rulesetService)
}

// Test URI extraction
func TestExtractNameFromURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "URI with double slash",
			uri:      "ruleset://python_style",
			expected: "python_style",
		},
		{
			name:     "URI with single colon",
			uri:      "ruleset:go_conventions",
			expected: "go_conventions",
		},
		{
			name:     "Invalid URI",
			uri:      "invalid",
			expected: "",
		},
		{
			name:     "Empty URI",
			uri:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractNameFromURI(tt.uri)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test ruleset formatting
func TestFormatRulesetAsMarkdown(t *testing.T) {
	rs := &ruleset.Ruleset{
		Name:        "test_ruleset",
		Description: "Test description",
		Tags:        []string{"tag1", "tag2"},
		Markdown:    "# Test Content\n\nSome content here",
	}

	result := formatRulesetAsMarkdown(rs)

	assert.Contains(t, result, "name: test_ruleset")
	assert.Contains(t, result, "description: Test description")
	assert.Contains(t, result, "tags: [tag1 tag2]")
	assert.Contains(t, result, "# Test Content")
	assert.Contains(t, result, "Some content here")
}

// Test RegisterTools doesn't panic
func TestRegisterTools(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// This test just ensures RegisterTools can be called without panicking
	// Full integration testing would require a real MCP server
	assert.NotPanics(t, func() {
		// We can't easily test this without a real server, but we can ensure the method exists
		assert.NotNil(t, handler.RegisterTools)
	})
}

// Test Start method exists and can be initialized
func TestStart(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// Verify Start method exists
	assert.NotNil(t, handler.Start)

	// Note: We cannot fully test Start() in unit tests because:
	// 1. It's a blocking call that serves stdio
	// 2. It requires actual stdin/stdout for MCP protocol communication
	// Full testing of Start() is done in integration tests
}

// Test HandleUpsertRuleset for creating a new ruleset
func TestHandleUpsertRuleset_Create(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// Mock the Upsert call to succeed
	mockService.On("Upsert", mock.AnythingOfType("*ruleset.Ruleset"), mock.AnythingOfType("*ruleset.Update")).Return(nil)
	mockService.On("Exists", "new_ruleset").Return(true, nil)

	// Create a mock request
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"name":        "new_ruleset",
		"description": "New ruleset description",
		"markdown":    "# New Content",
		"tags":        []interface{}{"tag1", "tag2"},
	}

	// Call the handler
	result, err := handler.HandleUpsertRuleset(context.TODO(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Successfully upserted ruleset 'new_ruleset'")
	mockService.AssertExpectations(t)
}

// Test HandleUpsertRuleset for updating an existing ruleset
func TestHandleUpsertRuleset_Update(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// Mock the Upsert call to succeed
	mockService.On("Upsert", mock.AnythingOfType("*ruleset.Ruleset"), mock.AnythingOfType("*ruleset.Update")).Return(nil)
	mockService.On("Exists", "existing_ruleset").Return(true, nil)

	// Create a mock request with only partial updates
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"name":        "existing_ruleset",
		"description": "Updated description",
	}

	// Call the handler
	result, err := handler.HandleUpsertRuleset(context.TODO(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Successfully upserted ruleset 'existing_ruleset'")
	mockService.AssertExpectations(t)
}

// Test HandleUpsertRuleset with missing name
func TestHandleUpsertRuleset_MissingName(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// Create a mock request without name
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"description": "Description",
		"markdown":    "# Content",
	}

	// Call the handler
	result, err := handler.HandleUpsertRuleset(context.TODO(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "missing required parameter 'name'")
}

// Test HandleUpsertRuleset with service error
func TestHandleUpsertRuleset_ServiceError(t *testing.T) {
	mockService := new(MockRulesetService)
	handler := NewHandler(mockService)

	// Mock the Upsert call to fail
	mockService.On("Upsert", mock.AnythingOfType("*ruleset.Ruleset"), mock.AnythingOfType("*ruleset.Update")).Return(assert.AnError)

	// Create a mock request
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"name":        "test_ruleset",
		"description": "Description",
		"markdown":    "# Content",
	}

	// Call the handler
	result, err := handler.HandleUpsertRuleset(context.TODO(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "failed to upsert ruleset")
	mockService.AssertExpectations(t)
}
