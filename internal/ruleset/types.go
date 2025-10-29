package ruleset

import "time"

// Ruleset represents a complete ruleset with all metadata and content
type Ruleset struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Tags         []string  `json:"tags"`
	Markdown     string    `json:"markdown"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}

// Update represents partial updates to an existing ruleset
type Update struct {
	Description *string   `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	Markdown    *string   `json:"markdown,omitempty"`
}
