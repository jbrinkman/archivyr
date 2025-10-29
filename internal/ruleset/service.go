package ruleset

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jbrinkman/archivyr/internal/util"
	"github.com/jbrinkman/archivyr/internal/valkey"
	"github.com/valkey-io/valkey-glide/go/v2/models"
)

// Service provides business logic for ruleset management
type Service struct {
	valkeyClient *valkey.Client
}

// NewService creates a new ruleset service instance
func NewService(client *valkey.Client) *Service {
	return &Service{
		valkeyClient: client,
	}
}

// Exists checks if a ruleset with the given name exists
func (s *Service) Exists(name string) (bool, error) {
	if err := util.ValidateRulesetName(name); err != nil {
		return false, err
	}

	key := fmt.Sprintf("ruleset:%s", name)
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	count, err := client.Exists(ctx, []string{key})
	if err != nil {
		return false, fmt.Errorf("failed to check if ruleset exists: %w", err)
	}

	return count > 0, nil
}

// ListNames retrieves all ruleset names from Valkey using SCAN
func (s *Service) ListNames() ([]string, error) {
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	names := make([]string, 0)
	cursor := models.NewCursor()

	// Use SCAN to iterate through all keys matching the pattern
	for {
		result, err := client.Scan(ctx, cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ruleset keys: %w", err)
		}

		// Extract names from keys that match the pattern (remove "ruleset:" prefix)
		for _, key := range result.Data {
			if len(key) > 8 && key[:8] == "ruleset:" { // len("ruleset:") = 8
				name := key[8:]
				names = append(names, name)
			}
		}

		cursor = result.Cursor
		if cursor.IsFinished() {
			break
		}
	}

	return names, nil
}

// Create creates a new ruleset in Valkey
func (s *Service) Create(ruleset *Ruleset) error {
	// Validate ruleset name
	if err := util.ValidateRulesetName(ruleset.Name); err != nil {
		return err
	}

	// Check if ruleset already exists
	exists, err := s.Exists(ruleset.Name)
	if err != nil {
		return err
	}

	if exists {
		// Get list of existing names for error message
		existingNames, listErr := s.ListNames()
		if listErr != nil {
			return fmt.Errorf("ruleset '%s' already exists", ruleset.Name)
		}
		return fmt.Errorf("ruleset '%s' already exists. Please choose a different name. Existing rulesets: %v", ruleset.Name, existingNames)
	}

	// Set timestamps
	now := time.Now()
	ruleset.CreatedAt = now
	ruleset.LastModified = now

	// Prepare hash fields
	key := fmt.Sprintf("ruleset:%s", ruleset.Name)
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	// Encode tags as JSON
	tagsJSON, err := json.Marshal(ruleset.Tags)
	if err != nil {
		return fmt.Errorf("failed to encode tags: %w", err)
	}

	// Store ruleset in Valkey hash
	fields := map[string]string{
		"description":   ruleset.Description,
		"tags":          string(tagsJSON),
		"markdown":      ruleset.Markdown,
		"created_at":    util.FormatTimestamp(ruleset.CreatedAt),
		"last_modified": util.FormatTimestamp(ruleset.LastModified),
	}

	_, err = client.HSet(ctx, key, fields)
	if err != nil {
		return fmt.Errorf("failed to create ruleset: %w", err)
	}

	return nil
}

// Get retrieves a ruleset by exact name from Valkey
func (s *Service) Get(name string) (*Ruleset, error) {
	// Validate ruleset name
	if err := util.ValidateRulesetName(name); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("ruleset:%s", name)
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	// Retrieve all hash fields
	result, err := client.HGetAll(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ruleset: %w", err)
	}

	// Check if ruleset exists (empty result means key doesn't exist)
	if len(result) == 0 {
		return nil, fmt.Errorf("ruleset '%s' not found", name)
	}

	// Parse hash fields into Ruleset struct
	ruleset := &Ruleset{
		Name: name,
	}

	// Extract and parse fields
	if desc, ok := result["description"]; ok {
		ruleset.Description = desc
	}

	if tagsJSON, ok := result["tags"]; ok {
		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}
		ruleset.Tags = tags
	}

	if markdown, ok := result["markdown"]; ok {
		ruleset.Markdown = markdown
	}

	if createdAtStr, ok := result["created_at"]; ok {
		createdAt, err := util.ParseTimestamp(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		ruleset.CreatedAt = createdAt
	}

	if lastModifiedStr, ok := result["last_modified"]; ok {
		lastModified, err := util.ParseTimestamp(lastModifiedStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse last_modified: %w", err)
		}
		ruleset.LastModified = lastModified
	}

	return ruleset, nil
}

// List retrieves all rulesets with metadata from Valkey
func (s *Service) List() ([]*Ruleset, error) {
	// Get all ruleset names
	names, err := s.ListNames()
	if err != nil {
		return nil, err
	}

	// Retrieve each ruleset
	rulesets := make([]*Ruleset, 0, len(names))
	for _, name := range names {
		ruleset, err := s.Get(name)
		if err != nil {
			// Skip rulesets that can't be retrieved (shouldn't happen, but be defensive)
			continue
		}
		rulesets = append(rulesets, ruleset)
	}

	return rulesets, nil
}

// Search searches for rulesets matching a glob pattern
func (s *Service) Search(pattern string) ([]*Ruleset, error) {
	if pattern == "" {
		return nil, fmt.Errorf("search pattern cannot be empty")
	}

	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	// Build the full key pattern for KEYS command
	keyPattern := fmt.Sprintf("ruleset:%s", pattern)

	// Use SCAN with pattern matching
	cursor := models.NewCursor()
	matchingNames := make([]string, 0)

	for {
		result, err := client.Scan(ctx, cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to search rulesets: %w", err)
		}

		// Filter keys that match our pattern and extract names
		for _, key := range result.Data {
			if len(key) > 8 && key[:8] == "ruleset:" {
				// Simple pattern matching - check if key matches the pattern
				if matchesPattern(key, keyPattern) {
					name := key[8:]
					matchingNames = append(matchingNames, name)
				}
			}
		}

		cursor = result.Cursor
		if cursor.IsFinished() {
			break
		}
	}

	// Retrieve full rulesets for matching names
	rulesets := make([]*Ruleset, 0, len(matchingNames))
	for _, name := range matchingNames {
		ruleset, err := s.Get(name)
		if err != nil {
			// Skip rulesets that can't be retrieved
			continue
		}
		rulesets = append(rulesets, ruleset)
	}

	return rulesets, nil
}

// Update updates an existing ruleset with the provided fields
func (s *Service) Update(name string, updates *RulesetUpdate) error {
	// Validate ruleset name
	if err := util.ValidateRulesetName(name); err != nil {
		return err
	}

	// Check if ruleset exists
	exists, err := s.Exists(name)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("ruleset '%s' not found", name)
	}

	// Prepare fields to update
	key := fmt.Sprintf("ruleset:%s", name)
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	fields := make(map[string]string)

	// Update only provided fields
	if updates.Description != nil {
		fields["description"] = *updates.Description
	}

	if updates.Tags != nil {
		tagsJSON, err := json.Marshal(*updates.Tags)
		if err != nil {
			return fmt.Errorf("failed to encode tags: %w", err)
		}
		fields["tags"] = string(tagsJSON)
	}

	if updates.Markdown != nil {
		fields["markdown"] = *updates.Markdown
	}

	// Always update last_modified timestamp
	fields["last_modified"] = util.FormatTimestamp(time.Now())

	// If no fields to update, return early
	if len(fields) == 1 { // Only last_modified
		return nil
	}

	// Update the hash in Valkey
	_, err = client.HSet(ctx, key, fields)
	if err != nil {
		return fmt.Errorf("failed to update ruleset: %w", err)
	}

	return nil
}

// Delete removes a ruleset from Valkey by name
func (s *Service) Delete(name string) error {
	// Validate ruleset name
	if err := util.ValidateRulesetName(name); err != nil {
		return err
	}

	// Check if ruleset exists
	exists, err := s.Exists(name)
	if err != nil {
		return err
	}

	if !exists {
		// Get list of existing names for error message
		existingNames, listErr := s.ListNames()
		if listErr != nil {
			return fmt.Errorf("ruleset '%s' not found", name)
		}
		return fmt.Errorf("ruleset '%s' not found. Existing rulesets: %v", name, existingNames)
	}

	// Delete the ruleset from Valkey
	key := fmt.Sprintf("ruleset:%s", name)
	ctx := s.valkeyClient.GetContext()
	client := s.valkeyClient.GetClient()

	_, err = client.Del(ctx, []string{key})
	if err != nil {
		return fmt.Errorf("failed to delete ruleset: %w", err)
	}

	return nil
}

// matchesPattern performs simple glob pattern matching
// Supports * (any characters) and ? (single character)
func matchesPattern(text, pattern string) bool {
	// Simple implementation for basic glob patterns
	// This is a basic version - for production, consider using filepath.Match or similar

	i, j := 0, 0
	starIdx, matchIdx := -1, 0

	for i < len(text) {
		if j < len(pattern) && (pattern[j] == '?' || pattern[j] == text[i]) {
			i++
			j++
		} else if j < len(pattern) && pattern[j] == '*' {
			starIdx = j
			matchIdx = i
			j++
		} else if starIdx != -1 {
			j = starIdx + 1
			matchIdx++
			i = matchIdx
		} else {
			return false
		}
	}

	for j < len(pattern) && pattern[j] == '*' {
		j++
	}

	return j == len(pattern)
}
