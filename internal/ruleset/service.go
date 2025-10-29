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
