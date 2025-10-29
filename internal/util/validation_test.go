package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRulesetName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		// Valid snake_case names
		{
			name:      "simple lowercase name",
			input:     "python",
			wantError: false,
		},
		{
			name:      "snake_case with underscores",
			input:     "python_style_guide",
			wantError: false,
		},
		{
			name:      "snake_case with numbers",
			input:     "python3_guide",
			wantError: false,
		},
		{
			name:      "snake_case with trailing number",
			input:     "style_guide_2",
			wantError: false,
		},
		{
			name:      "single letter",
			input:     "a",
			wantError: false,
		},
		{
			name:      "name with multiple underscores",
			input:     "my_python_style_guide",
			wantError: false,
		},
		// Invalid names
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "starts with uppercase",
			input:     "Python_guide",
			wantError: true,
		},
		{
			name:      "contains uppercase",
			input:     "python_Guide",
			wantError: true,
		},
		{
			name:      "starts with underscore",
			input:     "_python_guide",
			wantError: true,
		},
		{
			name:      "ends with underscore",
			input:     "python_guide_",
			wantError: true,
		},
		{
			name:      "double underscore",
			input:     "python__guide",
			wantError: true,
		},
		{
			name:      "contains spaces",
			input:     "python guide",
			wantError: true,
		},
		{
			name:      "contains hyphens",
			input:     "python-guide",
			wantError: true,
		},
		{
			name:      "contains special characters",
			input:     "python@guide",
			wantError: true,
		},
		{
			name:      "starts with number",
			input:     "3python_guide",
			wantError: true,
		},
		{
			name:      "camelCase",
			input:     "pythonGuide",
			wantError: true,
		},
		{
			name:      "PascalCase",
			input:     "PythonGuide",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRulesetName(tt.input)
			if tt.wantError {
				assert.Error(t, err, "expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "expected no error for input: %s", tt.input)
			}
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "standard timestamp",
			input:    time.Date(2025, 10, 28, 10, 30, 0, 0, time.UTC),
			expected: "2025-10-28T10:30:00Z",
		},
		{
			name:     "timestamp with timezone offset",
			input:    time.Date(2025, 10, 28, 10, 30, 0, 0, time.FixedZone("EST", -5*3600)),
			expected: "2025-10-28T10:30:00-05:00",
		},
		{
			name:     "timestamp with nanoseconds",
			input:    time.Date(2025, 10, 28, 10, 30, 0, 123456789, time.UTC),
			expected: "2025-10-28T10:30:00Z",
		},
		{
			name:     "zero timestamp",
			input:    time.Time{},
			expected: "0001-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimestamp(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  time.Time
	}{
		{
			name:      "valid RFC3339 timestamp",
			input:     "2025-10-28T10:30:00Z",
			wantError: false,
			expected:  time.Date(2025, 10, 28, 10, 30, 0, 0, time.UTC),
		},
		{
			name:      "valid RFC3339 with timezone offset",
			input:     "2025-10-28T10:30:00-05:00",
			wantError: false,
			expected:  time.Date(2025, 10, 28, 10, 30, 0, 0, time.FixedZone("", -5*3600)),
		},
		{
			name:      "valid RFC3339 with nanoseconds",
			input:     "2025-10-28T10:30:00.123456789Z",
			wantError: false,
			expected:  time.Date(2025, 10, 28, 10, 30, 0, 123456789, time.UTC),
		},
		{
			name:      "invalid format - missing timezone",
			input:     "2025-10-28T10:30:00",
			wantError: true,
		},
		{
			name:      "invalid format - wrong separator",
			input:     "2025-10-28 10:30:00Z",
			wantError: true,
		},
		{
			name:      "invalid format - incomplete date",
			input:     "2025-10-28",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid format - random string",
			input:     "not a timestamp",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimestamp(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.True(t, tt.expected.Equal(result), "expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFormatAndParseTimestamp_RoundTrip(t *testing.T) {
	// Test that formatting and parsing are inverse operations
	// Note: RFC3339 truncates to seconds, so we use a timestamp without nanoseconds
	original := time.Date(2025, 10, 28, 15, 45, 30, 0, time.UTC)

	formatted := FormatTimestamp(original)
	parsed, err := ParseTimestamp(formatted)

	require.NoError(t, err)
	assert.True(t, original.Equal(parsed), "round trip failed: expected %v, got %v", original, parsed)
}
