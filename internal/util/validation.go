package util

import (
	"fmt"
	"regexp"
	"time"
)

// snakeCaseRegex matches valid snake_case identifiers
var snakeCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`)

// ValidateRulesetName validates that a ruleset name follows snake_case convention
func ValidateRulesetName(name string) error {
	if name == "" {
		return fmt.Errorf("ruleset name cannot be empty")
	}

	if !snakeCaseRegex.MatchString(name) {
		return fmt.Errorf("ruleset name must be in snake_case format (lowercase letters, numbers, and underscores only, starting with a letter): %s", name)
	}

	return nil
}

// FormatTimestamp converts a time.Time to RFC3339 format string
func FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimestamp parses an RFC3339 format string to time.Time
func ParseTimestamp(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp format (expected RFC3339): %w", err)
	}
	return t, nil
}
