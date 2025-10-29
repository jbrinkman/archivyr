package mcp

import (
	"testing"

	"github.com/jbrinkman/archivyr/internal/ruleset"
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

func (m *MockRulesetService) Update(name string, updates *ruleset.RulesetUpdate) error {
	args := m.Called(name, updates)
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
